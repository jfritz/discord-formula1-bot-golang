package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Configuration
var debugMode = true
var localTimeZone = "Local" // "America/Chicago"

var RemoteCalendarUrl = "https://f1calendar.com/download/f1-calendar_p1_p2_p3_q_gp.ics"
var CalendarFilenameDefault = "formula.1.downloaded.ics"
var CalendarFilename = flag.String("calendarFile", CalendarFilenameDefault, "A fully qualified or relative path to the .ics file used for events.")
var WebhookConfFilenameDefault = "webhook_url.conf"
var WebhookConfFilename = flag.String("webhookFile", WebhookConfFilenameDefault, "A fully qualified or relative path to your webhook_url.conf file. Defaults to webhook_url.conf.")

func main() {
	findValidConfigFiles(CalendarFilename, WebhookConfFilename)

	// Attmept download file and use that instead
	var RemoteCalendarFileDownloadLocation = "formula.1.downloaded.ics"
	err := downloadFile(RemoteCalendarFileDownloadLocation, RemoteCalendarUrl)
	if err != nil {
		log.Printf("Couldn't download fresh calendar file: %s", err.Error())
	} else {
		log.Printf("Using downloaded file %s", RemoteCalendarFileDownloadLocation)
		*CalendarFilename = RemoteCalendarFileDownloadLocation
	}

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

		if len(events) != 0 {
			for _, event := range events {
				eventsStr += SummarizeEvent(event, localTimeZone)
			}

			webhook.SendMessage(prefix + eventsStr + suffix)
		}

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

func downloadFile(filepath string, url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			fmt.Printf("Error trying to close download stream %s: %s", err.Error(), url)
			return
		}
	}(response.Body)

	outfile, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer func(outfile *os.File) {
		err := outfile.Close()
		if err != nil {
			panic(fmt.Sprintf("Error trying to close local file %s: %s", err.Error(), url))
		}
	}(outfile)

	// Write the file
	_, err = io.Copy(outfile, response.Body)
	if err != nil {
		return err
	}

	return nil
}
