package errors

import (
	"net/http"
)

type DuplicateLedgerError string

// Error represents the error message
func (e DuplicateLedgerError) Error() string {
	return string(e)
}

// ErrCode represents the Xendit error code
func (e DuplicateLedgerError) ErrCode() string {
	return DuplicateLedger
}

// StatusCode represents the HTTP status code
func (e DuplicateLedgerError) StatusCode() int {
	return http.StatusConflict
}
