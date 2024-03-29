package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type DiscordWebhook struct {
	debug      bool
	configFile string
}

var WebhookHeaders = map[string]string{
	"Content-Type":              "application/x-www-form-urlencoded",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Encoding":           "gzip, deflate, br",
	"Accept-Language":           "en-US,en;q=0.5",
	"DNT":                       "1",
	"Host":                      "discordapp.com",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0",
}

func (webhook *DiscordWebhook) getWebhookUrl() string {
	file, err := os.Open(webhook.configFile)
	if err != nil {
		log.Print("need a valid config file")
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// By design only read the first webhook in the file
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Allow comment after the webhook, so throw it away here
	webhookUrl := strings.Split(scanner.Text(), " ")[0]

	// TODO validation on the webhook URL?
	return strings.TrimSpace(webhookUrl)
}

func (webhook *DiscordWebhook) SendMessage(msg string) {
	client := &http.Client{}
	form := url.Values{}
	form.Add("content", msg)
	req, err := http.NewRequest("POST", webhook.getWebhookUrl(), strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatalf("Couldn't create webhook request: %s", err.Error())
	}

	for key, value := range WebhookHeaders {
		req.Header.Add(key, value)
	}

	if webhook.debug {
		fmt.Println("Sending webhook request...")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if webhook.debug {
		fmt.Println("Webhook response: ")
		fmt.Println(resp.StatusCode)
		var body []byte
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		fmt.Println()
	}
}
