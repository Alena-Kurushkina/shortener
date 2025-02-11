package authenticator

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkAuthMiddleware(b *testing.B) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	authHandler := AuthMiddleware(nextHandler)

	req, err := http.NewRequest("POST", "http://test", nil)
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		authHandler.ServeHTTP(httptest.NewRecorder(), req)
	}
}
