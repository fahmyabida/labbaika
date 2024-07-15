package errors

import (
	"net/http"
)

type LedgerNotFoundError string

// Error represents the error message
func (e LedgerNotFoundError) Error() string {
	return string(e)
}

// ErrCode represents the Xendit error code
func (e LedgerNotFoundError) ErrCode() string {
	return LedgerNotFound
}

// StatusCode represents the HTTP status code
func (e LedgerNotFoundError) StatusCode() int {
	return http.StatusBadRequest
}
