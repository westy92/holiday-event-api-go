package main

import (
	"fmt"

	holidays "github.com/westy92/holiday-api-go"
)

func main() {
	client := holidays.New("<your API key>")

	events, err := client.GetEvents(holidays.GetEventsRequest{
		//Date:     "today",
		//Timezone: "America/Chicago",
		//Adult:    false,
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(events)
	}
}
