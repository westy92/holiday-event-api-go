package holidays

// An interface of the API's standard response.
type StandardResponseInterface interface {
}

// The API's standard response.
type StandardResponse struct {
	RateLimit RateLimit // The API plan's current rate limit and status
}

// Your API plan's current Rate Limit and status. Upgrade to increase these limits.
type RateLimit struct {
	Limit     int // The amount of requests allowed this billing cycle
	Remaining int // The amount of requests remaining this billing cycle
}

// The Request struct for calling GetEvents.
type GetEventsRequest struct {
	Date     string // Date to get the events for. Defaults to today.
	Adult    bool   // Include events that may be unsafe for viewing at work or by children. Default is false.
	Timezone string // IANA Time Zone for calculating dates and times. Defaults to America/Chicago.
}

// The Response struct returned by GetEvents.
type GetEventsResponse struct {
	StandardResponse                // Standard response fields
	Adult            bool           `json:"adult"`             // Whether Adult entries can be included
	Date             string         `json:"date"`              // The Date string
	Timezone         string         `json:"timezone"`          // The Timezone used to calculate the Date's Events
	Events           []EventSummary `json:"events"`            // The Date's Events
	MultidayStarting []EventSummary `json:"multiday_starting"` // Multi-day Events that start on Date
	MultidayOngoing  []EventSummary `json:"multiday_ongoing"`  // Multi-day Events that are continuing their observance on Date
}

// A summary of an Event.
type EventSummary struct {
	ID   string `json:"id"`   // The Event Id
	Name string `json:"name"` // The Event name
	URL  string `json:"url"`  // The Event URL
}

// The Request struct for calling Search.
type SearchRequest struct {
	Query string // The search query. Must be at least 3 characters long.
	Adult bool   // Include events that may be unsafe for viewing at work or by children. Default is false.
}

// The Response struct returned by Search.
type SearchResponse struct {
	StandardResponse                // Standard response fields
	Query            string         `json:"query"`  // The search query
	Adult            bool           `json:"adult"`  // Whether Adult entries can be included
	Events           []EventSummary `json:"events"` // The found Events
}

// The Request struct for calling GetEventInfo.
type GetEventInfoRequest struct {
	ID    string `json:"id"`    // The ID of the requested Event.
	Start int    `json:"start"` // The starting range of returned occurrences. Optional, defaults to 2 years prior.
	End   int    `json:"end"`   // The ending range of returned occurrences. Optional, defaults to 3 years in the future.
}

// The Response struct returned by GetEventInfo.
type GetEventInfoResponse struct {
	StandardResponse           // Standard response fields
	Event            EventInfo `json:"event"` // The Event Info
}

// Information about an Event's Pattern.
type Pattern struct {
	FirstYear        int    `json:"first_year"`        // The first year this event is observed (0 implies none or unknown)
	LastYear         int    `json:"last_year"`         // The last year this event is observed (0 implies none or unknown)
	Observed         string `json:"observed"`          // A description of how this event is observed (formatted as plain text)
	ObservedHTML     string `json:"observed_html"`     // A description of how this event is observed (formatted as HTML)
	ObservedMarkdown string `json:"observed_markdown"` // A description of how this event is observed (formatted as Markdown)
	Length           int    `json:"length"`            // For how many days this event is celebrated
}

// Information about an Event's Occurrence.
type Occurrence struct {
	Date   string `json:"date"`   // The date or timestamp the Event occurs
	Length int    `json:"length"` // The length (in days) of the Event occurrence
}

// Formatted Text.
type RichText struct {
	Text     string `json:"text"`     // Formatted as plain text
	HTML     string `json:"html"`     // Formatted as HTML
	Markdown string `json:"markdown"` // Formatted as Markdown
}

// Information about an Event's Alternate Name.
type AlternateName struct {
	Name      string `json:"name"`       // An Event's Alternate Name
	FirstYear int    `json:"first_year"` // The first year this Alternate Name was in effect (0 implies none or unknown)
	LastYear  int    `json:"last_year"`  // The last year this Alternate Name was in effect (0 implies none or unknown)
}

// Information about an Event image.
type ImageInfo struct {
	Small  string `json:"small"`  // A small image
	Medium string `json:"medium"` // A medium image
	Large  string `json:"large"`  // A large image
}

// Information about an Event Founder.
type FounderInfo struct {
	Name string `json:"name"` // The Founder's name
	URL  string `json:"url"`  // A link to the Founder
	Date string `json:"date"` // The date the Event was founded
}

// Information about an Event.
type EventInfo struct {
	EventSummary
	Adult          bool            `json:"adult"`           // Whether this Event is unsafe for children or viewing at work
	AlternateNames []AlternateName `json:"alternate_names"` // The Event's Alternate Names
	Hashtags       []string        `json:"hashtags"`        // The Event's hashtags
	Image          ImageInfo       `json:"image"`           // The Event's images
	Sources        []string        `json:"sources"`         // The Event's sources
	Description    RichText        `json:"description"`     // The Event's description
	HowToObserve   RichText        `json:"how_to_observe"`  // How to observe the Event
	Patterns       []Pattern       `json:"patterns"`        // Patterns defining when the Event is observed
	Occurrences    []Occurrence    `json:"occurrences"`     // The Event Occurrences (when it occurs)
	Founders       []FounderInfo   `json:"founders"`        // The Event's founders
}

// An Error response object.
type errorResponse struct {
	Error string `json:"error"` // A descriptive error message
}
