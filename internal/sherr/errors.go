// Package sherr defines errors for shortener service
package sherr

import "fmt"

// AlreadyExistError defines error in case of creating shortening for long URL that already exist in data storage
type AlreadyExistError struct {
	ExistShortStr    string
	ExistOriginalURL string
}

func (ex *AlreadyExistError) Error() string {
	return fmt.Sprintf("You are trying to shorten URL %s that already has shortening %s", ex.ExistOriginalURL, ex.ExistShortStr)
}

func NewAlreadyExistError(originalURL, shortURL string) error {
	return &AlreadyExistError{
		ExistShortStr:    shortURL,
		ExistOriginalURL: originalURL,
	}
}
