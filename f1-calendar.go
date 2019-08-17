package main

import (
	"fmt"
	"github.com/apognu/gocal"
	"log"
	"os"
	"time"
)

type F1Calendar struct {
	debug bool
	filename string
}

var FilterMap = map[string]time.Duration{
	"24h": time.Hour * 24,
	"7d": time.Hour * 24 *7,
}

func summarizeEvent(event gocal.Event) string {
	// TODO implement this
	desc := `:arrow_right: **%s**:
- Starts in: **%s**
- Starts at: **%s**
`
	return fmt.Sprintf(desc, event.Summary, event.Start, event.Location)
}

func (cal *F1Calendar) GetEvents(upcomingOnly bool, dateFilter string) []gocal.Event {
	f, err := os.Open(cal.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// TODO debug this
	//start, end := time.Now(), time.Now().Add(-1 * FilterMap[dateFilter])

	c := gocal.NewParser(f)
	//c.Start, c.End = &start, &end
	c.Parse()

	if cal.debug {
		for _, e := range c.Events {
			fmt.Printf(summarizeEvent(e))
		}
	}

	return c.Events
}

