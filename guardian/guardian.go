package guardian

import (
	"fmt"
	"log"
	"strconv"

	hotconfig "github.com/jarzamendia/hotconfig"
	"github.com/jarzamendia/swarmguardian/configuration"
	boltdb "github.com/jarzamendia/swarmguardian/data"
	"github.com/jarzamendia/swarmguardian/models"
	"github.com/jarzamendia/swarmguardian/notify"
)

//Verify x
func Verify(service models.Service, config configuration.GlobalConfig, maxValue string) {

	const offset = 3600

	currentTimeS := boltdb.GetTimestamp()

	currentTime, error := strconv.Atoi(currentTimeS)

	if error != nil {

		log.Fatal(error)

	}

	replicas, error := strconv.Atoi(service.Replicas)

	if error != nil {

		log.Fatal(error)

	}

	maxValueInt, error := strconv.Atoi(maxValue)

	if error != nil {

		log.Fatal(error)

	}

	debug := hotconfig.GetEnvVarOrDefault("DEBUG", "false")

	if debug == "true" {

		fmt.Println("currentTimeS: " + currentTimeS)

		fmt.Println("currentTime")
		fmt.Println(currentTime)

		fmt.Println("replicas")
		fmt.Println(replicas)
		fmt.Println("maxValueInt")
		fmt.Println(maxValueInt)
	}

	if replicas > maxValueInt {

		if debug == "true" {

			fmt.Println("The service have more than the ReplicasLimit.")
		}

		query := boltdb.Get("My", service.ServiceName)

		if query == "" {

			query = currentTimeS

			if debug == "true" {

				fmt.Println("The service has never pass the limit. So changing the Service Timestamp to:")
				fmt.Println(query)

			}

		} else {

			if debug == "true" {

				fmt.Println("This service sometime have passed the limit. His timestamp is:")
				fmt.Println(query)

			}

		}

		serviceTime, error := strconv.Atoi(query)

		if error != nil {

			log.Fatal(error)

		}

		if serviceTime == currentTime {

			errMsg := "The service " + service.ServiceName + " has been published with " + service.Replicas + " replicas."

			if debug == "true" {

				fmt.Println("serviceTime == currentTime.")

				log.Println(errMsg)

			}

			boltdb.Save("My", service.ServiceName, currentTimeS)

			if config.Slack == true {

				if debug == "true" {

					fmt.Println("Using the Slack notification channel.")

				}

				notify.SendSlackNotification(config.SlackWebhookURL, config.SlackUsername, errMsg)

			}

			if config.Stdout == true {

				if debug == "true" {

					fmt.Println("Using the Stdout notification channel.")

				}

				notify.SendStdoutNotification(errMsg)

			}

		} else if serviceTime > (currentTime + offset) {

			errMsg := "The service " + service.ServiceName + " has been published with " + service.Replicas + " replicas."

			if debug == "true" {

				fmt.Println("serviceTime > (currentTime + offset)")

				fmt.Println("serviceTime == currentTime.")

				log.Println(errMsg)
			}

			boltdb.Save("My", service.ServiceName, currentTimeS)

			if config.Slack == true {

				if debug == "true" {

					fmt.Println("Using the Slack notification channel.")

				}

				notify.SendSlackNotification(config.SlackWebhookURL, config.SlackUsername, errMsg)

			}

			if config.Stdout == true {

				if debug == "true" {

					fmt.Println("Using the Stdout notification channel.")

				}

				notify.SendStdoutNotification(errMsg)

			}

		} else {

			if debug == "true" {

				errMsg := "The service " + service.ServiceName + " has been published with more replicas than the limit, But we already informed this!"

				fmt.Println(errMsg)
			}

		}

	}

	return

}
