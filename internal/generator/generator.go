// Package generator realises function for generating random string of certain length
package generator

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateRandomString returns string of random characters of passed length.
func GenerateRandomString(length int) string {
	result := strings.Builder{}
	result.Grow(length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result.WriteByte(charset[n.Int64()])
	}

	return result.String()
}
