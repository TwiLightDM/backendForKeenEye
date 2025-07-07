package usecases

import (
	"errors"
)

var (
	AccountNotFoundError     = errors.New("account not found")
	UserAccountNotFoundError = errors.New("user account not found")
	DifferentPasswordError   = errors.New("passwords are not similar")
	HashPasswordError        = errors.New("failed to hash password")
	CreateError              = errors.New("failed to create entity")
	ReadError                = errors.New("failed to read entity")
	UpdateError              = errors.New("failed to update entity")
	DeleteError              = errors.New("failed to delete entity")
	NoFieldsError            = errors.New("no fields provided to update")
	MissingIdError           = errors.New("missing id field")
	ValidationError          = errors.New("validation failed")
)
