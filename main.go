package main

import (
	"time"
)

// Configuration
var debugMode = true
var CalendarFile = "formula.1.2019.ics"
var WebhookConfFilename = "webhook_url.conf"
var localTimeZone = "America/Chicago"

func main() {
	var webhook = DiscordWebhook{
		debug:      debugMode,
		configFile: WebhookConfFilename,
	}

	var cal = F1Calendar{
		debug:    debugMode,
		filename: CalendarFile,
	}

	dow := time.Now().Weekday()

	switch dow {

	case time.Monday:
		prefix := "<:f1:436383126743285760> Happy Monday! Here is the schedule for the next race weekend: \n"
		suffix := "<:nico:436342726309445643>"
		events := cal.GetEvents(Next7Days)
		eventsStr := ""

		for _, event := range events {
			eventsStr += SummarizeEvent(event, localTimeZone)
		}

		outputMessage := prefix + eventsStr + suffix
		if len(events) == 0 {
			eventsStr = "No events in the next 7 days.\n"
		}
		webhook.SendMessage(outputMessage)

	case time.Thursday, time.Friday, time.Saturday, time.Sunday:
		prefix := "<:f1:436383126743285760> Race Weekend! In the next 24 hours: \n"
		suffix := "<:nico:436342726309445643>"
		events := cal.GetEvents(Next24Hours)
		eventsStr := ""

		for _, event := range events {
			eventsStr += SummarizeEvent(event, localTimeZone)
		}

		outputMessage := prefix + eventsStr + suffix
		if len(events) == 0 {
			eventsStr = "No events in the next 24 hours.\n"
		}
		webhook.SendMessage(outputMessage)

		if len(events) > 0 && dow == time.Friday {
			reminder := "@F1 Fantasy League - It's Friday! Remember to check your F1 fantasy teams!"
			webhook.SendMessage(reminder)
		}
	default:
	}
}
