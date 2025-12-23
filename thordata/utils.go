package thordata

import (
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"
	"strings"
)

func BuildUserAgent(version string) string {
	return fmt.Sprintf("thordata-go-sdk/%s (go %s; %s/%s)", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func BuildAuthHeaders(scraperToken string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + scraperToken,
		"Content-Type":  "application/x-www-form-urlencoded",
	}
}

func BuildPublicHeaders(publicToken, publicKey string) map[string]string {
	return map[string]string{
		"token":        publicToken,
		"key":          publicKey,
		"Content-Type": "application/x-www-form-urlencoded",
	}
}

func ToFormBody(payload map[string]string) string {
	v := url.Values{}
	for k, val := range payload {
		if strings.TrimSpace(val) == "" {
			continue
		}
		v.Set(k, val)
	}
	return v.Encode()
}

func SafeParseJSON(data []byte) (any, error) {
	var obj any
	if err := json.Unmarshal(data, &obj); err == nil {
		// If API returns a JSON string that itself contains JSON, try one more pass.
		if s, ok := obj.(string); ok {
			var inner any
			if err2 := json.Unmarshal([]byte(strings.TrimSpace(s)), &inner); err2 == nil {
				return inner, nil
			}
		}
		return obj, nil
	}
	// Not JSON
	return string(data), nil
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		// Handles numbers (float64 from JSON), bools, etc.
		s := fmt.Sprint(v)
		if s == "<nil>" {
			return ""
		}
		return strings.TrimSpace(s)
	}
}
