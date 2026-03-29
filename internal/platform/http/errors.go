package platformhttp

import "errors"

var (
	ErrInvalidRequestBody = errors.New("http: invalid request body")
	ErrValidationFailed   = errors.New("http: validation failed")
)
