package thordata

import (
	"fmt"
)

type APIError struct {
	Message    string
	StatusCode int
	Code       int
	Payload    any
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s (status=%d, code=%d)", e.Message, e.StatusCode, e.Code)
}

type AuthError struct{ APIError }
type RateLimitError struct {
	APIError
	RetryAfterSeconds int
}
type ServerError struct{ APIError }
type ValidationError struct{ APIError }
type NotCollectedError struct{ APIError }

func RaiseForCode(message string, payload map[string]any, statusCode int) error {
	apiCode := 0
	if v, ok := payload["code"]; ok {
		if n, ok2 := v.(float64); ok2 {
			apiCode = int(n)
		}
	}
	// Precedence: payload code (when != 200) > HTTP status (when != 200)
	effective := 0
	if apiCode != 0 && apiCode != 200 {
		effective = apiCode
	} else if statusCode != 0 && statusCode != 200 {
		effective = statusCode
	} else {
		effective = apiCode
	}

	errMsg := message
	if v, ok := payload["msg"]; ok {
		errMsg = fmt.Sprintf("%v", v)
	} else if v, ok := payload["message"]; ok {
		errMsg = fmt.Sprintf("%v", v)
	}

	base := APIError{Message: errMsg, StatusCode: statusCode, Code: apiCode, Payload: payload}

	switch {
	case effective == 300:
		return &NotCollectedError{APIError: base}
	case effective == 401 || effective == 403:
		return &AuthError{APIError: base}
	case effective == 402 || effective == 429:
		return &RateLimitError{APIError: base}
	case effective >= 500 && effective < 600:
		return &ServerError{APIError: base}
	case effective == 400 || effective == 422:
		return &ValidationError{APIError: base}
	default:
		return &APIError{Message: errMsg, StatusCode: statusCode, Code: apiCode, Payload: payload}
	}
}