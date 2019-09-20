package configuration

//GlobalConfig A global configuraration.
type GlobalConfig struct {
	Slack              bool
	SlackServerAddress string
	SlackServerPort    string
	SlackWebhookURL    string
	SlackUsername      string
	SlackPassword      string
	Stdout             bool
	Mail               bool
}
