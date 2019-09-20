package notify

import (
	"fmt"
	"time"
)

//SendStdoutNotification Send a Slack notification
func SendStdoutNotification(msg string) {

	now := time.Now()
	fmt.Println(now.String() + " :: " + msg)

}
