package notify

import (
	slack "github.com/leprosus/golang-slack-notifier"
)

//SendSlackNotification Send a Slack notification
func SendSlackNotification(url string, username string, msg string) {

	slackConn := slack.New(url)

	// Send notification
	slackConn.Notify(username, msg)

}
