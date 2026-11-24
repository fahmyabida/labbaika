package errors

import (
	"net/http"
)

type CashbookNotFoundError string

// Error represents the error message
func (e CashbookNotFoundError) Error() string {
	return string(e)
}

// ErrCode represents the Xendit error code
func (e CashbookNotFoundError) ErrCode() string {
	return CashbookNotFound
}

// StatusCode represents the HTTP status code
func (e CashbookNotFoundError) StatusCode() int {
	return http.StatusBadRequest
}
