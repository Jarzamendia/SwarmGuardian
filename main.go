package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	hotconfig "github.com/jarzamendia/hotconfig"
	"github.com/jarzamendia/swarmguardian/configuration"
	"github.com/jarzamendia/swarmguardian/guardian"
	"github.com/jarzamendia/swarmguardian/models"
)

func main() {

	fmt.Println("...:::SwarmGuardian:::...")

	fmt.Println("Starting...")

	debug := hotconfig.GetEnvVarOrDefault("DEBUG", "false")

	if debug == "true" {

		fmt.Println("DUBUG active!")
	}

	usingSlack := hotconfig.GetEnvVarOrDefault("SLACK", "false")
	usingStdout := hotconfig.GetEnvVarOrDefault("STDOUT", "false")
	usingMail := hotconfig.GetEnvVarOrDefault("MAIL", "false")

	if usingSlack == "true" {

		fmt.Println("Using SLACK...")

	}

	if usingStdout == "true" {

		fmt.Println("Using STDOUT...")

	}

	if usingMail == "true" {

		fmt.Println("Using MAIL...")

	}

	if usingSlack == "false" && usingStdout == "false" && usingMail == "false" {

		fmt.Println("At least one notification channel must be configured.")

		log.Fatal("No Notification channels are configured.")

	}

	r := mux.NewRouter()

	r.HandleFunc("/reconfigure", reconfigureHandler).Methods("GET")

	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8081",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start Server
	go func() {

		log.Println("Starting Server...")

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)

}

func reconfigureHandler(w http.ResponseWriter, r *http.Request) {

	usingSlack := hotconfig.GetEnvVarOrDefault("SLACK", "false")
	usingStdout := hotconfig.GetEnvVarOrDefault("STDOUT", "false")
	usingMail := hotconfig.GetEnvVarOrDefault("MAIL", "false")
	debug := hotconfig.GetEnvVarOrDefault("DEBUG", "false")

	if debug == "true" {

		fmt.Println("Reconfigure called.")
	}

	config := configuration.GlobalConfig{
		Slack:  false,
		Stdout: false,
		Mail:   false,
	}

	// TODO Verificar se as variaveis estão vazias.

	if usingSlack == "true" {

		config.Slack = true

		config.SlackServerAddress = hotconfig.GetEnvVarOrDefault("SLACK_SERVERADDRES", "")
		config.SlackServerPort = hotconfig.GetEnvVarOrDefault("SLACK_SERVERPORT", "")
		config.SlackWebhookURL = hotconfig.GetEnvVarOrDefault("SLACK_WEBHOOKURL", "")
		config.SlackUsername = hotconfig.GetEnvVarOrDefault("SLACK_USERNAME", "")
		config.SlackPassword = hotconfig.GetEnvVarOrDefault("SLACK_PASSWORD", "")

	}

	if usingStdout == "true" {

		config.Stdout = true

	}

	if usingMail == "true" {

		config.Mail = true

	}

	var maxValue = hotconfig.GetEnvVarOrDefault("MAXREPLICAS", "3")

	service := models.Service{

		ServiceName: r.FormValue("serviceName"),
		Replicas:    r.FormValue("replicas"),
	}

	if debug == "true" {

		fmt.Println("ServiceName: " + string(service.ServiceName))
		fmt.Println("Replicas: " + string(service.Replicas))

	}

	guardian.Verify(service, config, maxValue)

	w.WriteHeader(http.StatusOK)

}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
