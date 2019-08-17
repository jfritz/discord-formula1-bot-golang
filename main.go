package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

// Configuration
var DoRequest bool = true
var RootDirectory string
var CalendarFile string = "./formula.1.2019.ics"
var WebhookConfFilename string = "webhook_url.conf"

func main() {
	var webhook = DiscordWebhook{
		debug: true,
		url:   getWebhookUrl(WebhookConfFilename),
	}

	// TODO instantiate calendar obj

	// 0 = Sunday, 1 = Monday, ..., 4 = Thursday, 5 = Friday, 6 = Saturday
	dow := time.Now().Weekday()

	switch dow {
	case time.Monday:
		// TODO get next race events and output
		webhook.SendMessage("This is a test message")
	case time.Thursday, time.Friday, time.Saturday, time.Sunday:
		// TODO get next 24h of events and output
		webhook.SendMessage("This is a test message")
	default:
	}
}

func getWebhookUrl(confFilename string) string {
	file, err := os.Open(confFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// By design only read the first webhook in the file
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(scanner.Text())
}
