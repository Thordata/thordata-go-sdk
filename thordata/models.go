// thordata/models.go (新建)
package thordata

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
