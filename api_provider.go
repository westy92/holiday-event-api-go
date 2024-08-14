package holidays

import (
	"net/http"
	"net/url"
	"strconv"
)

// TODO docs

type APIProvider int

const (
	APILayer APIProvider = iota
	RapidAPI
	// TODO APIMarket.
)

func (api APIProvider) isValid() bool {
	switch api {
	case APILayer:
		return true
	case RapidAPI:
		return true
	default:
		return false
	}
}

func (api APIProvider) apiKeySource() string {
	switch api {
	case APILayer:
		return "https://apilayer.com/marketplace/checkiday-api#pricing"
	case RapidAPI:
		return "https://rapidapi.com/westy92-llc-westy92-llc-default/api/checkiday/pricing"
	default:
		return ""
	}
}

func (api APIProvider) baseURL() url.URL {
	//nolint:golint,exhaustruct
	switch api {
	case APILayer:
		return url.URL{
			Scheme: "https",
			Host:   "api.apilayer.com",
			Path:   "checkiday",
		}
	case RapidAPI:
		return url.URL{
			Scheme: "https",
			Host:   "checkiday.p.rapidapi.com",
		}
	default:
		return url.URL{}
	}
}

func (api APIProvider) extractRateLimitInfo(headers http.Header) RateLimit {
	var limit, remaining int

	switch api {
	case APILayer:
		limit, _ = strconv.Atoi(headers.Get("X-Ratelimit-Limit-Month"))
		remaining, _ = strconv.Atoi(headers.Get("X-Ratelimit-Remaining-Month"))
	case RapidAPI:
		limit, _ = strconv.Atoi(headers.Get("X-Ratelimit-Requests-Limit"))
		remaining, _ = strconv.Atoi(headers.Get("X-Ratelimit-Requests-Remaining"))
	default:
	}

	return RateLimit{
		Limit:     limit,
		Remaining: remaining,
	}
}

func (api APIProvider) attachRequestHeaders(headers *http.Header, apiKey string) {
	switch api {
	case APILayer:
		headers.Set("Apikey", apiKey)
	case RapidAPI:
		headers.Set("X-Rapidapi-Key", apiKey)
		headers.Set("X-Rapidapi-Host", "checkiday.p.rapidapi.com")
	default:
	}
}
