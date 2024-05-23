package notify

import (
	"strings"
	"time"

	"github.com/ku9nov/backup/configs"
	"github.com/slack-go/slack"
)

func SendMessageToSlack(cfgValues configs.Config, dbName []string) {
	channelID := cfgValues.Slack.SlackChannelID
	dbNames := strings.Join(dbName, ", ")
	client := slack.New(cfgValues.Slack.SlackToken, slack.OptionDebug(false))
	attachment := slack.Attachment{
		Title: "Backup Notification",
		Text:  "Backup is failed",
		Color: "#FF0000",
		Fields: []slack.AttachmentField{
			{
				Title: "Host",
				Value: cfgValues.Default.Host,
				Short: true,
			},
			{
				Title: "Failed backup",
				Value: dbNames,
				Short: true,
			},
			{
				Title: "Event time",
				Value: time.Now().Format("2006.01.02 15:04:05"),
				Short: true,
			},
		},
	}
	_, timestamp, err := client.PostMessage(
		channelID,
		slack.MsgOptionAttachments(attachment),
	)
	_ = timestamp
	if err != nil {
		panic(err)
	}

}
