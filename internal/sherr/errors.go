// Package sherr defines errors for shortener service.
package sherr

import (
	"errors"
	"fmt"
)

// AlreadyExistError defines error in case of creating shortening for long URL that already exist in data storage.
type AlreadyExistError struct {
	ExistShortStr    string
	ExistOriginalURL string
}

// Error gives string representation of error in output.
func (ex *AlreadyExistError) Error() string {
	return fmt.Sprintf("You are trying to shorten URL %s that already has shortening %s", ex.ExistOriginalURL, ex.ExistShortStr)
}

// NewAlreadyExistError creates AlreadyExistError.
func NewAlreadyExistError(originalURL, shortURL string) error {
	return &AlreadyExistError{
		ExistShortStr:    shortURL,
		ExistOriginalURL: originalURL,
	}
}

// ErrNoUserIDInToken defines error in case of empty user ID in JWT.
var ErrNoUserIDInToken = errors.New("no user ID in JWT")

// ErrTokenInvalid defines error in case of invalid JWT.
var ErrTokenInvalid = errors.New("token is not valid")

// ErrDBRecordDeleted defines error in case of requesting deleted shortening.
var ErrDBRecordDeleted = errors.New("shortening is deleted")
