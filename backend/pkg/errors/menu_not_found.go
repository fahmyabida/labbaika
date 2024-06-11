package errors

import (
	"net/http"
)

type MenuNotFoundError string

// Error represents the error message
func (e MenuNotFoundError) Error() string {
	return string(e)
}

// ErrCode represents the Xendit error code
func (e MenuNotFoundError) ErrCode() string {
	return MenuNotFound
}

// StatusCode represents the HTTP status code
func (e MenuNotFoundError) StatusCode() int {
	return http.StatusBadRequest
}
