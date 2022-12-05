package main

import (
	"fmt"
	"math/rand"
	"time"

	holidays "github.com/westy92/holiday-event-api-go"
)

func main() {
	client, err := holidays.New("<your API key>")

	if err != nil {
		fmt.Println(err)
		return
	}

	events, err := client.GetEvents(holidays.GetEventsRequest{
		// These parameters are the defaults but can be specified:
		// Date:     "today",
		// Timezone: "America/Chicago",
		// Adult:    false,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	rand.Seed(time.Now().Unix())
	randomEvent := events.Events[rand.Intn(len(events.Events))]
	fmt.Printf("Today is %s! Find more information at: %s.\n", randomEvent.Name, randomEvent.Url)
	fmt.Printf("Rate limit remaining: %d/%d (month).\n", events.RateLimit.RemainingMonth, events.RateLimit.LimitMonth)
}
