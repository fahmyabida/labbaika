package errors

import (
	"net/http"
)

type DuplicateCategoryError string

// Error represents the error message
func (e DuplicateCategoryError) Error() string {
	return string(e)
}

// ErrCode represents the Xendit error code
func (e DuplicateCategoryError) ErrCode() string {
	return DuplicateCategory
}

// StatusCode represents the HTTP status code
func (e DuplicateCategoryError) StatusCode() int {
	return http.StatusConflict
}
