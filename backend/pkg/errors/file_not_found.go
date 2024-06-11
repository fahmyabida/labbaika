package errors

import (
	"net/http"
)

type FileNotFoundError string

// Error represents the error message
func (e FileNotFoundError) Error() string {
	return string(e)
}

// ErrCode represents the Xendit error code
func (e FileNotFoundError) ErrCode() string {
	return FileNotFound
}

// StatusCode represents the HTTP status code
func (e FileNotFoundError) StatusCode() int {
	return http.StatusBadRequest
}
