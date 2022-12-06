package main

import (
	"fmt"

	holidays "github.com/westy92/holiday-event-api-go"
)

func main() {
	// Get a FREE API key from https://apilayer.com/marketplace/checkiday-api#pricing
	client, err := holidays.New("<your API key>")

	if err != nil {
		fmt.Println(err)
		return
	}

	// Get Events for a given Date
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

	event := events.Events[0]
	fmt.Printf("Today is %s! Find more information at: %s.\n", event.Name, event.Url)
	fmt.Printf("Rate limit remaining: %d/%d (month).\n", events.RateLimit.RemainingMonth, events.RateLimit.LimitMonth)

	// Get Event Information
	eventInfo, err := client.GetEventInfo(holidays.GetEventInfoRequest{
		Id: event.Id,
		// These parameters can be specified to calculate the range of eventInfo.Event.Occurrences
		// Start: 2020,
		// End: 2030,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The Event's hashtags are %q.\n", eventInfo.Event.Hashtags)

	// Search for Events
	query := "pizza day"
	search, err := client.Search(holidays.SearchRequest{
		Query: query,
		// These parameters are the defaults but can be specified:
		// Adult: false
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Found %d events, including %s, that match the query \"%s\".\n", len(search.Events), search.Events[0].Name, query)
}
