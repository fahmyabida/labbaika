package errors

import (
	"net/http"
)

type CategoryNotFoundError string

// Error represents the error message
func (e CategoryNotFoundError) Error() string {
	return string(e)
}

// ErrCode represents the Xendit error code
func (e CategoryNotFoundError) ErrCode() string {
	return CategoryNotFound
}

// StatusCode represents the HTTP status code
func (e CategoryNotFoundError) StatusCode() int {
	return http.StatusBadRequest
}
