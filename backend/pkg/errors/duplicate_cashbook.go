package errors

import (
	"net/http"
)

type DuplicateCashbookError string

// Error represents the error message
func (e DuplicateCashbookError) Error() string {
	return string(e)
}

// ErrCode represents the Xendit error code
func (e DuplicateCashbookError) ErrCode() string {
	return DuplicateCashbook
}

// StatusCode represents the HTTP status code
func (e DuplicateCashbookError) StatusCode() int {
	return http.StatusConflict
}
