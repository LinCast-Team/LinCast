package entities

import "errors"

// Errors related with the user entity
var (
	ErrInvalidUsernameLength = errors.New("username must be between 3 and 50 characters")
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters")
	ErrPasswordRequirements  = errors.New("password must contain at least one number, one uppercase letter and one special character")
)
