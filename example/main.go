package main

import (
	"context"
	"fmt"

	holidays "github.com/westy92/holiday-event-api-go"
)

func main() {
	// Get a FREE API key from https://apilayer.com/marketplace/checkiday-api#pricing
	client, err := holidays.New(holidays.APILayer, "<your API key>")
	if err != nil {
		fmt.Println(err)

		return
	}

	ctx := context.TODO()

	// Get Events for a given Date
	events, err := client.GetEvents(ctx, holidays.GetEventsRequest{
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
	fmt.Printf("Today is %s! Find more information at: %s.\n", event.Name, event.URL)
	fmt.Printf("Rate limit remaining: %d/%d (billing cycle).\n", events.RateLimit.Remaining, events.RateLimit.Limit)

	// Get Event Information
	eventInfo, err := client.GetEventInfo(ctx, holidays.GetEventInfoRequest{
		ID: event.ID,
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

	search, err := client.Search(ctx, holidays.SearchRequest{
		Query: query,
		// These parameters are the defaults but can be specified:
		// Adult: false,
	})
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Printf("Found %d events, including %s, that match the query \"%s\".\n", len(search.Events), search.Events[0].Name, query)
}
