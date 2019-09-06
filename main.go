package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Configuration
var debugMode = true
var localTimeZone = "Local" // "America/Chicago"

var CalendarFilenameDefault = "formula.1.2019.ics"
var CalendarFilename = flag.String("calendarFile", CalendarFilenameDefault, "A fully qualified or relative path to the .ics file used for events.")
var WebhookConfFilenameDefault = "webhook_url.conf"
var WebhookConfFilename = flag.String("webhookFile", WebhookConfFilenameDefault, "A fully qualified or relative path to your webhook_url.conf file. Defaults to webhook_url.conf.")

func main() {
	findValidConfigFiles(CalendarFilename, WebhookConfFilename)

	var webhook = DiscordWebhook{
		debug:      debugMode,
		configFile: *WebhookConfFilename,
	}

	var cal = F1Calendar{
		debug:    debugMode,
		filename: *CalendarFilename,
	}

	if debugMode {
		fmt.Println("Using Calendar File: ", *CalendarFilename)
		fmt.Println("Using Webhook File: ", *WebhookConfFilename)
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

		if len(events) != 0 {
			for _, event := range events {
				eventsStr += SummarizeEvent(event, localTimeZone)
			}

			webhook.SendMessage(prefix + eventsStr + suffix)
		}

		if len(events) > 0 && dow == time.Friday {
			reminder := "<@&616676115645202462> - It's Friday! Remember to check your F1 fantasy teams!"
			webhook.SendMessage(reminder)
		}

	default:
	}
}

func findValidConfigFiles(calendarFile *string, webhookFile *string) {
	// Calendar File

	// Try the file as-given
	if !isValidFile(*calendarFile) {

		// Try to find calendar file in executable directory with given filepath
		localDir, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		localDir = filepath.Dir(localDir) // Trim off .exe

		candidateFile := localDir + string(os.PathSeparator) + *calendarFile
		if !isValidFile(candidateFile) {

			// Try to find calendar file in executable directory with default filename
			candidateFile = localDir + string(os.PathSeparator) + CalendarFilenameDefault
			if !isValidFile(candidateFile) {
				log.Fatal("Was not able to open the provided calendar file nor the default calendar file.")
			} else {
				*calendarFile = candidateFile
			}
		} else {
			*calendarFile = candidateFile
		}
	}

	// Configuration File

	// Try the file as-given
	if !isValidFile(*webhookFile) {

		// Try to find webhook file in executable directory with given filepath
		localDir, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		localDir = filepath.Dir(localDir) // Trim off .exe

		candidateFile := localDir + string(os.PathSeparator) + *webhookFile
		if !isValidFile(candidateFile) {

			// Try to find webhook file in executable directory with default filename
			candidateFile = localDir + string(os.PathSeparator) + WebhookConfFilenameDefault
			if !isValidFile(candidateFile) {
				log.Fatal("Was not able to open the provided webhook file nor the default webhook file.")
			} else {
				*webhookFile = candidateFile
			}
		} else {
			*webhookFile = candidateFile
		}
	}
}

func isValidFile(filename string) bool {
	ret := true
	if _, err := os.Open(filename); err != nil {
		ret = false
	}
	return ret
}
