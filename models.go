package holidays

type StandardResponseInterface interface {
}

type StandardResponse struct {
	RateLimit RateLimit
}

type RateLimit struct {
	LimitMonth     int
	RemainingMonth int
}

type GetEventsRequest struct {
	/**
	 * Date to get the events for. Defaults to today.
	 */
	Date string
	/**
	* Include events that may be unsafe for viewing at work or by children. Default is false.
	 */
	Adult bool
	/**
	* IANA Time Zone for calculating dates and times. Defaults to America/Chicago.
	 */
	Timezone string
}

type GetEventsResponse struct {
	StandardResponse
	Adult            bool           `json:"adult"`
	Date             string         `json:"date"` // TODO int or string!
	Timezone         string         `json:"timezone"`
	Events           []EventSummary `json:"events"`
	MultidayStarting []EventSummary `json:"multiday_starting"`
	MultidayOngoing  []EventSummary `json:"multiday_ongoing"`
}

type EventSummary struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type SearchRequest struct {
	/**
	 * The search query. Must be at least 3 characters long.
	 */
	Query string
	/**
	 * Include events that may be unsafe for viewing at work or by children. Default is false.
	 */
	Adult bool
}

type SearchResponse struct {
	StandardResponse
	Query  string         `json:"query"`
	Adult  bool           `json:"adult"`
	Events []EventSummary `json:"events"`
}

type GetEventInfoRequest struct {
	/**
	 * The ID of the requested Event.
	 */
	Id string `json:"id"`
	/**
	 * The starting range of returned occurrences. Optional, defaults to 2 years prior.
	 */
	Start int `json:"start"`
	/**
	 * The ending range of returned occurrences. Optional, defaults to 3 years in the future.
	 */
	End int `json:"end"`
}

type GetEventInfoResponse struct {
	StandardResponse
	Event EventInfo `json:"event"`
}

type Pattern struct {
	FirstYear        int    `json:"first_year"`
	LastYear         int    `json:"last_year"`
	Observed         string `json:"observed"`
	ObservedHtml     string `json:"observed_html"`
	ObservedMarkdown string `json:"observed_markdown"`
	Length           int    `json:"length"`
}

type Occurrence struct {
	Date   string `json:"date"`
	Length int    `json:"length"`
}

type RichText struct {
	Text     string `json:"text"`
	Html     string `json:"html"`
	Markdown string `json:"markdown"`
}

type AlternateName struct {
	Name      string `json:"name"`
	FirstYear int    `json:"first_year"`
	LastYear  int    `json:"last_year"`
}

type ImageInfo struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type EventInfo struct {
	EventSummary
	Adult          bool            `json:"adult"`
	AlternateNames []AlternateName `json:"alternate_names"`
	Hashtags       []string        `json:"hashtags"`
	Image          ImageInfo       `json:"image"`
	Sources        []string        `json:"sources"`
	Description    RichText        `json:"description"`
	HowToObserve   RichText        `json:"how_to_observe"`
	Patterns       []Pattern       `json:"patterns"`
	Occurrences    []Occurrence    `json:"occurrences"`
}

type errorResponse struct {
	Error string `json:"error"`
}
