# The Official Holiday and Event API for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/westy92/holiday-event-api-go.svg)](https://pkg.go.dev/github.com/westy92/holiday-event-api-go)
[![Build Status](https://github.com/westy92/holiday-event-api-go/actions/workflows/github-actions.yml/badge.svg)](https://github.com/westy92/holiday-event-api-go/actions)
[![Code Coverage](https://codecov.io/gh/westy92/holiday-event-api-go/branch/main/graph/badge.svg)](https://codecov.io/gh/westy92/holiday-event-api-go)
[![Funding Status](https://img.shields.io/github/sponsors/westy92)](https://github.com/sponsors/westy92)

Industry-leading Holiday and Event API for Go. Over 5,000 holidays and thousands of descriptions. Trusted by the World’s leading companies. Built by developers for developers since 2011.

## Supported Go Versions

Latest version of the the Holiday and Event API supports last two Go major [releases](https://go.dev/doc/devel/release#policy) and might work with older versions.

## Authentication

Access to the Holiday and Event API requires an API Key. You can get for one for FREE [here](https://apilayer.com/marketplace/checkiday-api#pricing), no credit card required! Note that free plans are limited. To access more data and have more requests, a paid plan is required.

## Installation

```console
go get github.com/westy92/holiday-event-api-go
```

## Example

```go
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
```
