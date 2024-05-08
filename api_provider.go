package holidays

import (
	"net/http"
	"strconv"
)

// TODO docs

type ApiProvider int

const (
	ApiLayer ApiProvider = iota
	RapidApi
)

func (api ApiProvider) isValid() bool {
	switch api {
	case ApiLayer:
		return true
	case RapidApi:
		return true
	default:
		return false
	}
}

func (api ApiProvider) apiKeySource() string {
	switch api {
	case ApiLayer:
		return "https://apilayer.com/marketplace/checkiday-api#pricing"
	case RapidApi:
		return "https://rapidapi.com/westy92-llc-westy92-llc-default/api/checkiday/pricing"
	default:
		return ""
	}
}

func (api ApiProvider) baseUrl() string {
	switch api {
	case ApiLayer:
		return "https://api.apilayer.com/checkiday/"
	case RapidApi:
		return "https://checkiday.p.rapidapi.com/"
	default:
		return ""
	}
}

func (api ApiProvider) extractRateLimitInfo(headers http.Header) RateLimit {
	var limit, remaining int
	switch api {
	case ApiLayer:
		limit, _ = strconv.Atoi(headers.Get("x-ratelimit-limit-month"))
		remaining, _ = strconv.Atoi(headers.Get("x-ratelimit-remaining-month"))
	case RapidApi:
		limit, _ = strconv.Atoi(headers.Get("x-ratelimit-requests-limit"))
		remaining, _ = strconv.Atoi(headers.Get("x-ratelimit-requests-remaining"))
	default:
	}
	return RateLimit{
		Limit:     limit,
		Remaining: remaining,
	}
}

func (api ApiProvider) attachRequestHeaders(headers *http.Header, apiKey string) {
	switch api {
	case ApiLayer:
		headers.Set("apikey", apiKey)
	case RapidApi:
		headers.Set("X-RapidAPI-Key", apiKey)
		headers.Set("X-RapidAPI-Host", "checkiday.p.rapidapi.com")
	default:
	}
}
