package errors

import (
	"net/http"
)

type DuplicateMenuError string

// Error represents the error message
func (e DuplicateMenuError) Error() string {
	return string(e)
}

// ErrCode represents the Xendit error code
func (e DuplicateMenuError) ErrCode() string {
	return DuplicateMenu
}

// StatusCode represents the HTTP status code
func (e DuplicateMenuError) StatusCode() int {
	return http.StatusConflict
}
