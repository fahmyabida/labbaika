package errors

import (
	"net/http"
)

// InvalidInputError is a generic 400 for request validation failures — the
// message is the specific reason (e.g. "amount must be greater than zero"),
// not a fixed per-field error type, since validation messages vary freely.
type InvalidInputError string

// Error represents the error message
func (e InvalidInputError) Error() string {
	return string(e)
}

// ErrCode represents the Xendit error code
func (e InvalidInputError) ErrCode() string {
	return InvalidInput
}

// StatusCode represents the HTTP status code
func (e InvalidInputError) StatusCode() int {
	return http.StatusBadRequest
}
