package holidays

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// The API Client
type Client struct {
	apiKey string
}

const (
	version   = "0.0.2"
	userAgent = "HolidayApiGo/" + version
	baseUrl   = "https://api.apilayer.com/checkiday/"
)

// Creates a New Client using the provided API key.
// Get a FREE API key from https://apilayer.com/marketplace/checkiday-api#pricing
func New(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("please provide a valid API key. Get one at https://apilayer.com/marketplace/checkiday-api#pricing")
	}

	return &Client{
		apiKey: apiKey,
	}, nil
}

// Gets the Events for the provided Date
func (c *Client) GetEvents(req GetEventsRequest) (*GetEventsResponse, error) {
	var params = url.Values{
		"adult": {strconv.FormatBool(req.Adult)},
	}

	if req.Timezone != "" {
		params["timezone"] = []string{req.Timezone}
	}

	if req.Date != "" {
		params["date"] = []string{req.Date}
	}

	res, rateLimit, err := request[GetEventsResponse](c, "events", params)
	if err != nil {
		return nil, err
	}

	res.RateLimit = *rateLimit

	return res, nil
}

// Gets the Event Info for the provided Event
func (c *Client) GetEventInfo(req GetEventInfoRequest) (*GetEventInfoResponse, error) {
	var params = url.Values{}

	if req.Id == "" {
		return nil, errors.New("event id is required")
	}
	params["id"] = []string{req.Id}

	if req.Start != 0 {
		params["start"] = []string{strconv.Itoa(req.Start)}
	}

	if req.End != 0 {
		params["end"] = []string{strconv.Itoa(req.End)}
	}

	res, rateLimit, err := request[GetEventInfoResponse](c, "event", params)
	if err != nil {
		return nil, err
	}

	res.RateLimit = *rateLimit

	return res, nil
}

// Searches for Events with the given criteria
func (c *Client) Search(req SearchRequest) (*SearchResponse, error) {
	var params = url.Values{
		"adult": {strconv.FormatBool(req.Adult)},
	}

	if req.Query == "" {
		return nil, errors.New("search query is required")
	}
	params["query"] = []string{req.Query}

	res, rateLimit, err := request[SearchResponse](c, "search", params)
	if err != nil {
		return nil, err
	}

	res.RateLimit = *rateLimit

	return res, nil
}

// Gets the API Client Version
func (c *Client) GetVersion() string {
	return version
}

func request[R StandardResponseInterface](client *Client, urlPath string, params url.Values) (*R, *RateLimit, error) {
	url, err := url.Parse(baseUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("can't parse baseUrl: %w", err)
	}
	url.Path = path.Join(url.Path, urlPath)

	if params != nil {
		url.RawQuery = params.Encode()
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create request: %w", err)
	}

	req.Header.Set("apikey", client.apiKey)
	req.Header.Set("User-Agent", userAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("can't process request: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		var errBody errorResponse
		if err := json.NewDecoder(res.Body).Decode(&errBody); err == nil && errBody.Error != "" {
			return nil, nil, errors.New(errBody.Error)
		}
		return nil, nil, errors.New(res.Status)
	}

	var result R
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("can't parse response: %w", err)
	}

	limitMonth, _ := strconv.Atoi(res.Header.Get("x-ratelimit-limit-month"))
	remainingMonth, _ := strconv.Atoi(res.Header.Get("x-ratelimit-remaining-month"))
	rateLimit := RateLimit{
		LimitMonth:     limitMonth,
		RemainingMonth: remainingMonth,
	}

	return &result, &rateLimit, nil
}
