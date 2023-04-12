package sdk

import "errors"

var (
	ErrPrincipalIRNIsNotSet = errors.New("principal IRN is not set")
	ErrAPIKeyIsEmpty        = errors.New("API key is empty")
)
