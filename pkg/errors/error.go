package errors

import "net/http"

type RestError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

func NewRestError(statusCode int, message string) *RestError {
	return &RestError{
		Message:    message,
		StatusCode: statusCode,
		Error:      http.StatusText(statusCode),
	}
}

func MalformedJWTError(message string) *RestError {
	return NewRestError(http.StatusBadRequest, message) // returns RestError
}

func InternalServerError(message string) *RestError {
	return NewRestError(http.StatusInternalServerError, message) // returns RestError
}
