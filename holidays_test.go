package holidays_test

import (
	"context"
	"errors"
	"runtime"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	holidays "github.com/westy92/holiday-event-api-go"
)

var ErrTest = errors.New("err")

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("fails with missing API Key", func(t *testing.T) {
		t.Parallel()

		api, err := holidays.New(holidays.APILayer, "")

		assert := assert.New(t)
		assert.Nil(api)
		require.EqualError(t, err, "please provide a valid API key. Get one at https://apilayer.com/marketplace/checkiday-api#pricing")
	})

	t.Run("fails with invalid provider", func(t *testing.T) {
		t.Parallel()

		api, err := holidays.New(holidays.APIProvider(-1), "abc123")

		assert := assert.New(t)
		assert.Nil(api)
		require.EqualError(t, err, "please provide a valid API provider")
	})

	t.Run("returns a new client", func(t *testing.T) {
		t.Parallel()

		api, err := holidays.New(holidays.APILayer, "abc123")

		assert := assert.New(t)
		require.NoError(t, err)
		assert.NotNil(api)
	})
}

func TestCommonFunctionality(t *testing.T) {
	t.Run("passes along API key", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			MatchHeader("apikey", "abc123").
			Reply(200).
			File("testdata/getEvents-default.json")

		api, _ := holidays.New(holidays.APILayer, "abc123")
		_, _ = api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		assert.True(gock.IsDone())
	})

	t.Run("passes along user-agent", func(t *testing.T) {
		api, _ := holidays.New(holidays.APILayer, "abc123")

		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			MatchHeader("user-agent", "HolidayApiGo/"+api.GetVersion()).
			Reply(200).
			File("testdata/getEvents-default.json")

		_, _ = api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		assert.True(gock.IsDone())
	})

	t.Run("passes along platform version", func(t *testing.T) {
		api, _ := holidays.New(holidays.APILayer, "abc123")

		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			MatchHeader("X-Platform-Version", runtime.Version()).
			Reply(200).
			File("testdata/getEvents-default.json")

		_, _ = api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		assert.True(gock.IsDone())
	})

	t.Run("passes along error", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(401).
			JSON(map[string]string{"error": "MyError!"})

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		assert.Nil(response)
		require.Error(t, err, "MyError!")

		assert.True(gock.IsDone())
	})

	t.Run("server error (500)", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(500)

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		assert.Nil(response)
		require.EqualError(t, err, "error status returned from API: 500 Internal Server Error")

		assert.True(gock.IsDone())
	})

	t.Run("server error (unknown)", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(599)

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		assert.Nil(response)
		require.EqualError(t, err, "error status returned from API: 599 ")

		assert.True(gock.IsDone())
	})

	t.Run("server error (other)", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			ReplyError(ErrTest)

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		assert.Nil(response)
		require.Error(t, err, "can't process request: Get \"https://api.apilayer.com/checkiday/events?adult=false\": err")

		assert.True(gock.IsDone())
	})

	t.Run("server error (malformed response)", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(200).
			JSON("{")

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		assert.Nil(response)
		require.EqualError(t, err, "can't parse response: unexpected EOF")

		assert.True(gock.IsDone())
	})

	t.Run("follows redirects", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(302).
			SetHeader("Location", "https://api.apilayer.com/checkiday/redirected")

		gock.New("https://api.apilayer.com/checkiday/").
			Get("/redirected").
			Reply(200).
			File("testdata/getEvents-default.json")

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		require.NoError(t, err)
		assert.Equal("America/Chicago", response.Timezone)

		assert.True(gock.IsDone())
	})

	t.Run("reports rate limits", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(200).
			SetHeader("X-RateLimit-Limit-Month", "100").
			SetHeader("x-ratelimit-remaining-month", "88").
			File("testdata/getEvents-default.json")

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		require.NoError(t, err)
		assert.Equal(holidays.RateLimit{
			Limit:     100,
			Remaining: 88,
		}, response.RateLimit)

		assert.True(gock.IsDone())
	})

	t.Run("nil context", func(t *testing.T) {
		api, _ := holidays.New(holidays.APILayer, "abc123")
		//nolint:golint,staticcheck
		response, err := api.GetEvents(nil, holidays.GetEventsRequest{})

		assert := assert.New(t)
		assert.Nil(response)
		require.Error(t, err)
		assert.Equal("can't create request: net/http: nil Context", err.Error())
	})
}

func TestGetEvents(t *testing.T) {
	t.Run("fetches with default parameters", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/events").
			Reply(200).
			File("testdata/getEvents-default.json")

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEvents(context.TODO(), holidays.GetEventsRequest{})

		assert := assert.New(t)
		require.NoError(t, err)
		assert.False(response.Adult)
		assert.Equal("America/Chicago", response.Timezone)
		assert.Len(response.Events, 2)
		assert.Len(response.MultidayStarting, 1)
		assert.Len(response.MultidayOngoing, 2)
		assert.Equal(holidays.EventSummary{
			ID:   "b80630ae75c35f34c0526173dd999cfc",
			Name: "Cinco de Mayo",
			URL:  "https://www.checkiday.com/b80630ae75c35f34c0526173dd999cfc/cinco-de-mayo",
		}, response.Events[0])

		assert.True(gock.IsDone())
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

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEvents(context.TODO(), holidays.GetEventsRequest{
			Adult:    true,
			Timezone: "America/New_York",
			Date:     "7/16/1992",
		})

		assert := assert.New(t)
		require.NoError(t, err)
		assert.True(response.Adult)
		assert.Equal("America/New_York", response.Timezone)
		assert.Len(response.Events, 2)
		assert.Empty(response.MultidayStarting)
		assert.Len(response.MultidayOngoing, 1)
		assert.Equal(holidays.EventSummary{
			ID:   "6ebb6fd5e483de2fde33969a6c398472",
			Name: "Get to Know Your Customers Day",
			URL:  "https://www.checkiday.com/6ebb6fd5e483de2fde33969a6c398472/get-to-know-your-customers-day",
		}, response.Events[0])

		assert.True(gock.IsDone())
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

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEventInfo(context.TODO(), holidays.GetEventInfoRequest{
			ID: "f90b893ea04939d7456f30c54f68d7b4",
		})

		assert := assert.New(t)
		require.NoError(t, err)
		assert.Equal("f90b893ea04939d7456f30c54f68d7b4", response.Event.ID)
		assert.Len(response.Event.Hashtags, 2)

		assert.True(gock.IsDone())
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

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEventInfo(context.TODO(), holidays.GetEventInfoRequest{
			ID:    "f90b893ea04939d7456f30c54f68d7b4",
			Start: 2002,
			End:   2003,
		})

		assert := assert.New(t)
		require.NoError(t, err)
		assert.Len(response.Event.Occurrences, 2)
		assert.Equal(holidays.Occurrence{
			Date:   "08/08/2002",
			Length: 1,
		}, response.Event.Occurrences[0])

		assert.True(gock.IsDone())
	})

	t.Run("invalid event", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/event").
			MatchParam("id", "hi").
			Reply(404).
			JSON(map[string]string{"error": "Event not found."})

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEventInfo(context.TODO(), holidays.GetEventInfoRequest{
			ID: "hi",
		})

		assert := assert.New(t)
		assert.Nil(response)
		require.EqualError(t, err, "error returned from API: Event not found.")

		assert.True(gock.IsDone())
	})

	t.Run("missing id", func(t *testing.T) {
		t.Parallel()

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.GetEventInfo(context.TODO(), holidays.GetEventInfoRequest{})

		assert := assert.New(t)
		assert.Nil(response)
		require.EqualError(t, err, "event id is required")
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

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.Search(context.TODO(), holidays.SearchRequest{
			Query: "zucchini",
		})

		assert := assert.New(t)
		require.NoError(t, err)
		assert.False(response.Adult)
		assert.Equal("zucchini", response.Query)
		assert.Len(response.Events, 3)
		assert.Equal(holidays.EventSummary{
			ID:   "cc81cbd8730098456f85f69798cbc867",
			Name: "National Zucchini Bread Day",
			URL:  "https://www.checkiday.com/cc81cbd8730098456f85f69798cbc867/national-zucchini-bread-day",
		}, response.Events[0])

		assert.True(gock.IsDone())
	})

	t.Run("fetches with set parameters", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/search").
			MatchParam("query", "porch day").
			MatchParam("adult", "true").
			Reply(200).
			File("testdata/search-parameters.json")

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.Search(context.TODO(), holidays.SearchRequest{
			Query: "porch day",
			Adult: true,
		})

		assert := assert.New(t)
		require.NoError(t, err)
		assert.True(response.Adult)
		assert.Equal("porch day", response.Query)
		assert.Len(response.Events, 1)
		assert.Equal(holidays.EventSummary{
			ID:   "61363236f06e4eb8e4e14e5925c2503d",
			Name: "Sneak Some Zucchini Onto Your Neighbor's Porch Day",
			URL:  "https://www.checkiday.com/61363236f06e4eb8e4e14e5925c2503d/sneak-some-zucchini-onto-your-neighbors-porch-day",
		}, response.Events[0])

		assert.True(gock.IsDone())
	})

	t.Run("query too short", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/search").
			MatchParam("query", "a").
			Reply(400).
			JSON(map[string]string{"error": "Please enter a longer search term."})

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.Search(context.TODO(), holidays.SearchRequest{
			Query: "a",
		})

		assert := assert.New(t)
		assert.Nil(response)
		require.EqualError(t, err, "error returned from API: Please enter a longer search term.")

		assert.True(gock.IsDone())
	})

	t.Run("too many results", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.apilayer.com/checkiday/").
			Get("/search").
			MatchParam("query", "day").
			Reply(400).
			JSON(map[string]string{"error": "Too many results returned. Please refine your query."})

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.Search(context.TODO(), holidays.SearchRequest{
			Query: "day",
		})

		assert := assert.New(t)
		assert.Nil(response)
		require.EqualError(t, err, "error returned from API: Too many results returned. Please refine your query.")

		assert.True(gock.IsDone())
	})

	t.Run("missing parameters", func(t *testing.T) {
		t.Parallel()

		api, _ := holidays.New(holidays.APILayer, "abc123")
		response, err := api.Search(context.TODO(), holidays.SearchRequest{})

		assert := assert.New(t)
		assert.Nil(response)
		require.EqualError(t, err, "search query is required")
	})
}
