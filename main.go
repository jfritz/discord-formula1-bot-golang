package main

import (
	"time"
)

// Configuration
var	debugMode bool = true
var CalendarFile string = "formula.1.2019.ics"
var WebhookConfFilename string = "webhook_url.conf"
var localTimeZone = "America/Chicago"

func main() {
	var webhook = DiscordWebhook{
		debug: 		debugMode,
		configFile: WebhookConfFilename,
	}

	var cal = F1Calendar{
		debug:    debugMode,
		filename: CalendarFile,
	}

	dow := time.Now().Weekday()

	switch dow {
	case time.Monday:
		events := cal.GetEvents(Next7Days)
		// TODO use SummarizeEvent(e, tz) and build message to send
		webhook.SendMessage("This is a test message")
	case time.Thursday, time.Friday, time.Saturday, time.Sunday:
		events := cal.GetEvents(Next24Hours)
		// TODO use SummarizeEvent(e, tz) and build message to send
		webhook.SendMessage("This is a test message")
		if dow == time.Friday {
			// TODO also reminder to update fantasy F1 league
		}
	default:
	}
}
