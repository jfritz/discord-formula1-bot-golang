package main

import (
	"fmt"
	"github.com/apognu/gocal"
	"log"
	"os"
	"strings"
	"time"
)

type F1Calendar struct {
	debug    bool
	filename string
}

const Next24Hours time.Duration = time.Hour * 24
const Next7Days time.Duration = time.Hour * 24 * 7

// Returns a human readable string to describe the length of the given duration
// Assumes you don't care about the number of seconds remaining
func friendlyDurationString(duration time.Duration) string {
	// Thank you to https://freshman.tech/golang-timer/
	total := int(duration.Seconds())
	days := int(total / (60 * 60 * 24))
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60

	daysStr, hoursStr, minutesStr := "", "", ""
	if days > 0 {
		daysStr = fmt.Sprintf("%d days,", days)
	}
	if hours > 0 {
		hoursStr = fmt.Sprintf("%d hours,", hours)
	}
	if minutes > 0 {
		minutesStr = fmt.Sprintf("%d minutes", minutes)
	}

	return strings.TrimSpace(strings.Join([]string{daysStr, hoursStr, minutesStr}, " "))
}

// Returns a human readable string, summarizing the F1 event provided in `event`
// Time is adjusted to the given `localTimeZone`
func SummarizeEvent(event gocal.Event, localTimeZone string) string {
	localTime, err := time.LoadLocation(localTimeZone)
	if err != nil {
		log.Fatal(err)
	}

	desc := `:arrow_right: **%s** @ %s: 
- Starts in **%s** at **%s**

`
	return fmt.Sprintf(desc,
		event.Summary,
		event.Location,
		friendlyDurationString(time.Until(event.Start.UTC())),
		event.Start.In(localTime))
}

// Returns a slice of `gocal.Event`s between now and the given duration
func (cal *F1Calendar) GetEvents(filter time.Duration) []gocal.Event {
	f, err := os.Open(cal.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	start, end := time.Now(), time.Now().Add(filter)

	c := gocal.NewParser(f)
	c.Start, c.End = &start, &end
	c.Parse()

	if cal.debug {
		fmt.Println("Found events: ")
		for _, e := range c.Events {
			fmt.Printf(SummarizeEvent(e, "America/Chicago"))
		}
	}

	return c.Events
}
