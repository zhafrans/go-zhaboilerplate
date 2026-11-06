package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given Param is not valid")
	// ErrUnauthorized will throw if the credentials are not valid or unauthorized
	ErrUnauthorized = errors.New("unauthorized: invalid credentials or token")
	// ErrForbidden will throw if the user is forbidden or inactive
	ErrForbidden = errors.New("forbidden: account is inactive or access denied")
)
