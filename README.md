[![Go Reference](https://pkg.go.dev/badge/github.com/westy92/holiday-event-api-go.svg)](https://pkg.go.dev/github.com/westy92/holiday-event-api-go)
[![Build Status](https://github.com/westy92/holiday-event-api-go/actions/workflows/github-actions.yml/badge.svg)](https://github.com/westy92/holiday-event-api-go/actions)
[![Code Coverage](https://codecov.io/gh/westy92/holiday-event-api-go/branch/main/graph/badge.svg)](https://codecov.io/gh/westy92/holiday-event-api-go)
[![Known Vulnerabilities](https://snyk.io/test/github/westy92/holiday-event-api-go/badge.svg)](https://snyk.io/test/github/westy92/holiday-event-api-go)
[![Funding Status](https://img.shields.io/github/sponsors/westy92)](https://github.com/sponsors/westy92)

Industry-leading Holiday and Event API for JavaScript/TypeScript. Over 5,000 holidays and thousands of descriptions. Trusted by the World’s leading companies. Built by developers for developers since 2011.

## Supported Go Versions
Latest version of the the Holiday and Event API supports last two Go major [releases](https://go.dev/doc/devel/release#policy) and might work with older versions.

## Authentication

Access to the Holiday and Event API requires an API Key. You can get for one for FREE [here](https://apilayer.com/marketplace/checkiday-api#pricing), no credit card required! Note that free plans are limited. To access more data and have more requests, a paid plan is required.

## Installation

```
go get westy92/holiday-event-api-go/v0.0.1
```

## Example

```go
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
```