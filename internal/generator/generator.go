// Package generator realises function for generating random string of certain length
package generator

import (
	"strings"
	"time"

	"math/rand"
)

// GenerateRandomString returns string of random characters of passed length.
func GenerateRandomString(length int) string {
	charset := []byte{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := strings.Builder{}
	result.Grow(length)
	for i := 0; i < length; i++ {
		result.WriteByte(charset[random.Intn(len(charset))])
	}

	return result.String()
}
