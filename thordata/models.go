// thordata/models.go
package thordata

// Standard API Response wrapper
type APIResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
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
