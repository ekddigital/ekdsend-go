package ekdsend

import "fmt"

// EKDSendError is the base error type for API errors
type EKDSendError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	RequestID  string `json:"request_id"`
}

func (e *EKDSendError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("EKDSend API error: %s (code: %s, status: %d, request_id: %s)",
			e.Message, e.Code, e.StatusCode, e.RequestID)
	}
	return fmt.Sprintf("EKDSend API error: %s (code: %s, status: %d)",
		e.Message, e.Code, e.StatusCode)
}

// AuthenticationError is returned when API key is invalid (401)
type AuthenticationError struct {
	EKDSendError
}

// ValidationError is returned when request validation fails (400)
type ValidationError struct {
	EKDSendError
	Errors map[string]interface{} `json:"errors"`
}

// RateLimitError is returned when rate limit is exceeded (429)
type RateLimitError struct {
	EKDSendError
	RetryAfter int `json:"retry_after"`
}

// NotFoundError is returned when resource is not found (404)
type NotFoundError struct {
	EKDSendError
}

// IsAuthenticationError checks if the error is an authentication error
func IsAuthenticationError(err error) bool {
	_, ok := err.(*AuthenticationError)
	return ok
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsRateLimitError checks if the error is a rate limit error
func IsRateLimitError(err error) bool {
	_, ok := err.(*RateLimitError)
	return ok
}

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
