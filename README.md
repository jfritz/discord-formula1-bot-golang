# discord-formula1-bot-golang

WARNING: This bot is not yet complete.

A simple discord bot to post F1 events using webhooks

This script is designed to be run once every 24 hours.

Mondays -
    Send a message with the full next race weekend
    
Thursdays, Fridays, Saturdays, Sundays -
    Send a message with events that will occur in the next 24h (if any)


Note: This is based off of my Python project, discord-formula1-bot, commit 1eb0c68.

## Requirements
- Golang system requirements for Golang 1.12 apply: https://github.com/golang/go/wiki/MinimumRequirements

## Installation
- Place your webhook url as the first line in webhook_url.conf in the same folder as the binary.
