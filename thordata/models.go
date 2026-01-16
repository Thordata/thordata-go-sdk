// thordata/models.go
package thordata

import "time"

// --- Common ---

// APIResponse is the standard JSON wrapper {"code":..., "msg":..., "data":...}
type APIResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
	// Some endpoints put data at root level alongside code, handled by custom logic or specific structs
}

// --- Scraper / SERP ---

type SerpResponse struct {
	Code           int            `json:"code"`
	Organic        []SearchResult `json:"organic"`
	SearchMetadata SearchMetadata `json:"search_metadata"`
	// Add other fields as needed (news, shopping, etc.)
	RawJSON map[string]any `json:"-"` // Fallback for extra fields
}

type SearchResult struct {
	Title         string `json:"title"`
	Link          string `json:"link"`
	Snippet       string `json:"snippet"`
	Position      int    `json:"position"`
	DisplayedLink string `json:"displayed_link"`
}

type SearchMetadata struct {
	Id             string  `json:"id"`
	Status         string  `json:"status"`
	JsonEndpoint   string  `json:"json_endpoint"`
	CreatedAt      string  `json:"created_at"`
	ProcessedAt    string  `json:"processed_at"`
	TotalTimeTaken float64 `json:"total_time_taken"`
}

// --- Universal ---

type UniversalResponse struct {
	Code int    `json:"code"`
	HTML string `json:"html,omitempty"`
	PNG  string `json:"png,omitempty"` // Base64
}

// --- Tasks ---

type TaskCreateResponse struct {
	TaskId string `json:"task_id"`
}

type TaskStatus struct {
	TaskId string `json:"task_id"`
	Status string `json:"status"`
	// Add other fields if API returns more
}

type TaskDownload struct {
	Download string `json:"download"`
}

type TaskList struct {
	Count int          `json:"count"`
	List  []TaskStatus `json:"list"`
}

// --- Public / Account ---

type UsageStatistics struct {
	TotalUsageTraffic float64      `json:"total_usage_traffic"`
	TrafficBalance    float64      `json:"traffic_balance"`
	QueryDays         int          `json:"query_days"`
	RangeUsageTraffic float64      `json:"range_usage_traffic"`
	Data              []UsageDaily `json:"data"`
}

type UsageDaily struct {
	Date         string  `json:"date"`
	UsageTraffic float64 `json:"usage_traffic"`
}

// --- Proxy Users ---

type ProxyUser struct {
	Username     string  `json:"username"`
	Status       any     `json:"status"` // API might return bool or int or string
	TrafficLimit int64   `json:"traffic_limit"`
	UsageTraffic float64 `json:"usage_traffic"`
}

type ProxyUserList struct {
	Limit          float64     `json:"limit"`
	RemainingLimit float64     `json:"remaining_limit"`
	UserCount      int         `json:"user_count"`
	List           []ProxyUser `json:"list"`
}

// --- Proxy List ---

type ProxyServer struct {
	IP             string `json:"ip"`
	Port           int    `json:"port"`
	Username       string `json:"username"`        // Sometimes 'user'
	Password       string `json:"password"`        // Sometimes 'pwd'
	ExpirationTime any    `json:"expiration_time"` // int or string
	Region         string `json:"region"`
}

// --- Locations ---

type Country struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
}

type State struct {
	StateCode string `json:"state_code"`
	StateName string `json:"state_name"`
}

type City struct {
	CityName string `json:"city_name"`
}

type ASN struct {
	AsnCode string `json:"asn_code"`
	AsnName string `json:"asn_name"`
}

// CommonSettings for video/audio tasks
type CommonSettings struct {
	Resolution        string `json:"resolution,omitempty"`
	AudioFormat       string `json:"audio_format,omitempty"`
	Bitrate           string `json:"bitrate,omitempty"`
	IsSubtitles       string `json:"is_subtitles,omitempty"`
	SubtitlesLanguage string `json:"subtitles_language,omitempty"`
}

type VideoTaskOptions struct {
	FileName       string
	SpiderID       string
	SpiderName     string
	Parameters     map[string]any
	CommonSettings CommonSettings
	IncludeErrors  bool
}

// RunTaskConfig configures the polling behavior for RunTask.
type RunTaskConfig struct {
	// MaxWait is the maximum time to wait for task completion. Default: 10 minutes.
	MaxWait time.Duration
	// InitialPollInterval is the starting interval between status checks. Default: 2 seconds.
	InitialPollInterval time.Duration
	// MaxPollInterval is the cap for the exponential backoff. Default: 10 seconds.
	MaxPollInterval time.Duration
}
