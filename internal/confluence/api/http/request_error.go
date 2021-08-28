package http

import (
	"fmt"
)

// RequestError encapsulates an HTTP request error, complete with status code.
type RequestError struct {
	StatusCode   int
	ResponseBody string
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("HTTP response error: %d %s", r.StatusCode, r.ResponseBody)
}

// newRequestError creates an HTTP request error, complete with status code.
func newRequestError(statusCode int, responseBody string) error {
	return &RequestError{
		StatusCode:   statusCode,
		ResponseBody: responseBody,
	}
}
