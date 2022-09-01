package holidays

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestCommonFunctionality(t *testing.T) {
	t.Run("passes along API key", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			MatchHeader("apikey", "abc123").
			Reply(200).
			File("testdata/getEvents-default.json")

		api := New("abc123")
		api.GetEvents(GetEventsRequest{})

		assert.True(t, gock.IsDone())
	})

	t.Run("passes along user-agent", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			MatchHeader("user-agent", "HolidayApiGo/0.0.1"). // TODO grab version from somewhere?
			Reply(200).
			File("testdata/getEvents-default.json")

		api := New("abc123")
		api.GetEvents(GetEventsRequest{})

		assert.True(t, gock.IsDone())
	})

	t.Run("passes along error", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(401).
			JSON(map[string]string{"error": "MyError!"})

		api := New("abc123")
		response, err := api.GetEvents(GetEventsRequest{})

		assert.Nil(t, response)
		assert.EqualError(t, err, "MyError!")

		assert.True(t, gock.IsDone())
	})

	t.Run("server error", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(500)

		api := New("abc123")
		response, err := api.GetEvents(GetEventsRequest{})

		assert.Nil(t, response)
		assert.EqualError(t, err, "500 Internal Server Error")

		assert.True(t, gock.IsDone())
	})

	t.Run("reports rate limits", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(200).
			SetHeader("X-RateLimit-Limit-Month", "100").
			SetHeader("x-ratelimit-remaining-month", "88").
			SetHeader("x-ratelimit-limit-day", "10").
			SetHeader("x-ratelimit-remaining-day", "9").
			File("testdata/getEvents-default.json")

		api := New("abc123")
		response, err := api.GetEvents(GetEventsRequest{})

		assert.Nil(t, err)
		assert.Equal(t, response.RateLimit, RateLimit{
			LimitMonth:     100,
			LimitDay:       10,
			RemainingMonth: 88,
			RemainingDay:   9,
		})

		assert.True(t, gock.IsDone())
	})
}

func TestGetEvents(t *testing.T) {
	t.Run("fetches with default parameters", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(200).
			File("testdata/getEvents-default.json")

		api := New("abc123")
		response, err := api.GetEvents(GetEventsRequest{})

		assert.Nil(t, err)
		assert.False(t, response.Adult)
		assert.Equal(t, response.Timezone, "America/Chicago")
		assert.Len(t, response.Events, 2)
		assert.Len(t, response.MultidayStarting, 1)
		assert.Len(t, response.MultidayOngoing, 2)
		assert.Equal(t, response.Events[0], EventSummary{
			Id:   "b80630ae75c35f34c0526173dd999cfc",
			Name: "Cinco de Mayo",
			Url:  "https://www.checkiday.com/b80630ae75c35f34c0526173dd999cfc/cinco-de-mayo",
		})

		assert.True(t, gock.IsDone())
	})

	t.Run("fetches with set parameters", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			MatchParam("adult", "true").
			MatchParam("timezone", "America/New_York").
			MatchParam("date", "7/16/1992").
			Reply(200).
			File("testdata/getEvents-parameters.json")

		api := New("abc123")
		response, err := api.GetEvents(GetEventsRequest{
			Adult:    true,
			Timezone: "America/New_York",
			Date:     "7/16/1992",
		})

		assert.Nil(t, err)
		assert.True(t, response.Adult)
		assert.Equal(t, response.Timezone, "America/New_York")
		assert.Len(t, response.Events, 2)
		assert.Len(t, response.MultidayStarting, 0)
		assert.Len(t, response.MultidayOngoing, 1)
		assert.Equal(t, response.Events[0], EventSummary{
			Id:   "6ebb6fd5e483de2fde33969a6c398472",
			Name: "Get to Know Your Customers Day",
			Url:  "https://www.checkiday.com/6ebb6fd5e483de2fde33969a6c398472/get-to-know-your-customers-day",
		})

		assert.True(t, gock.IsDone())
	})
}

func TestGetEventInfo(t *testing.T) {
	t.Run("fetches with default parameters", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/event").
			MatchParam("id", "f90b893ea04939d7456f30c54f68d7b4").
			Reply(200).
			File("testdata/getEventInfo.json")

		api := New("abc123")
		response, err := api.GetEventInfo(GetEventInfoRequest{
			Id: "f90b893ea04939d7456f30c54f68d7b4",
		})

		assert.Nil(t, err)
		assert.Equal(t, response.Event.Id, "f90b893ea04939d7456f30c54f68d7b4")
		assert.Len(t, response.Event.Hashtags, 2)

		assert.True(t, gock.IsDone())
	})

	t.Run("fetches with set parameters", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/event").
			MatchParam("id", "f90b893ea04939d7456f30c54f68d7b4").
			MatchParam("start", "2002").
			MatchParam("end", "2003").
			Reply(200).
			File("testdata/getEventInfo-parameters.json")

		api := New("abc123")
		response, err := api.GetEventInfo(GetEventInfoRequest{
			Id:    "f90b893ea04939d7456f30c54f68d7b4",
			Start: 2002,
			End:   2003,
		})

		assert.Nil(t, err)
		assert.Len(t, response.Event.Occurrences, 2)
		assert.Equal(t, response.Event.Occurrences[0], Occurrence{
			Date:   "08/08/2002",
			Length: 1,
		})

		assert.True(t, gock.IsDone())
	})

	t.Run("invalid event", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/event").
			MatchParam("id", "hi").
			Reply(404).
			JSON(map[string]string{"error": "Event not found."})

		api := New("abc123")
		response, err := api.GetEventInfo(GetEventInfoRequest{
			Id: "hi",
		})

		assert.Nil(t, response)
		assert.EqualError(t, err, "Event not found.")

		assert.True(t, gock.IsDone())
	})

	t.Run("missing id", func(t *testing.T) {
		api := New("abc123")
		response, err := api.GetEventInfo(GetEventInfoRequest{})

		assert.Nil(t, response)
		assert.EqualError(t, err, "Event id is required.")
	})
}

func TestSearch(t *testing.T) {
	t.Run("fetches with default parameters", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/search").
			MatchParam("query", "zucchini").
			Reply(200).
			File("testdata/search-default.json")

		api := New("abc123")
		response, err := api.Search(SearchRequest{
			Query: "zucchini",
		})

		assert.Nil(t, err)
		assert.False(t, response.Adult)
		assert.Equal(t, response.Query, "zucchini")
		assert.Len(t, response.Events, 3)
		assert.Equal(t, response.Events[0], EventSummary{
			Id:   "cc81cbd8730098456f85f69798cbc867",
			Name: "National Zucchini Bread Day",
			Url:  "https://www.checkiday.com/cc81cbd8730098456f85f69798cbc867/national-zucchini-bread-day",
		})

		assert.True(t, gock.IsDone())
	})

	t.Run("fetches with set parameters", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/search").
			MatchParam("query", "porch day").
			MatchParam("adult", "true").
			Reply(200).
			File("testdata/search-parameters.json")

		api := New("abc123")
		response, err := api.Search(SearchRequest{
			Query: "porch day",
			Adult: true,
		})

		assert.Nil(t, err)
		assert.True(t, response.Adult)
		assert.Equal(t, response.Query, "porch day")
		assert.Len(t, response.Events, 1)
		assert.Equal(t, response.Events[0], EventSummary{
			Id:   "61363236f06e4eb8e4e14e5925c2503d",
			Name: "Sneak Some Zucchini Onto Your Neighbor's Porch Day",
			Url:  "https://www.checkiday.com/61363236f06e4eb8e4e14e5925c2503d/sneak-some-zucchini-onto-your-neighbors-porch-day",
		})

		assert.True(t, gock.IsDone())
	})

	t.Run("query too short", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/search").
			MatchParam("query", "a").
			Reply(400).
			JSON(map[string]string{"error": "Please enter a longer search term."})

		api := New("abc123")
		response, err := api.Search(SearchRequest{
			Query: "a",
		})

		assert.Nil(t, response)
		assert.EqualError(t, err, "Please enter a longer search term.")

		assert.True(t, gock.IsDone())
	})

	t.Run("too many results", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/search").
			MatchParam("query", "day").
			Reply(400).
			JSON(map[string]string{"error": "Too many results returned. Please refine your query."})

		api := New("abc123")
		response, err := api.Search(SearchRequest{
			Query: "day",
		})

		assert.Nil(t, response)
		assert.EqualError(t, err, "Too many results returned. Please refine your query.")

		assert.True(t, gock.IsDone())
	})

	t.Run("missing parameters", func(t *testing.T) {
		api := New("abc123")
		response, err := api.Search(SearchRequest{})

		assert.Nil(t, response)
		assert.EqualError(t, err, "Search query is required.")
	})
}
