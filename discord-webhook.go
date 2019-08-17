package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type DiscordWebhook struct {
	debug bool
	url string
}

var WebhookHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
	"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Encoding": "gzip, deflate, br",
	"Accept-Language": "en-US,en;q=0.5",
	"DNT": "1",
	"Host": "discordapp.com",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0",
}

func (webhook *DiscordWebhook) SendMessage(msg string) {
	client := &http.Client{}
	form := url.Values{}
	form.Add("content", msg)
	req, err := http.NewRequest("POST", webhook.url, strings.NewReader(form.Encode()))

	for key,value := range WebhookHeaders {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	if webhook.debug {
		fmt.Println("Discord webhook response: ")
		fmt.Println(resp.StatusCode)
		var body []byte
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		fmt.Println()
	}
}



