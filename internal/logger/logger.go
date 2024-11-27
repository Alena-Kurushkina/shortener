package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Log represents global var for logging
// By default Log is no-op logger
var Log *zap.SugaredLogger = zap.NewNop().Sugar()

func Initialize() error {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"/Users/alena/log/shortener.log",
	}
	zl, err := cfg.Build()

	// zl, err := zap.NewProduction()
	if err != nil {
		return err
	}

	sugar := zl.Sugar()
	Log = sugar

	return nil
}

func logRequest(uri, method string, duration time.Duration) {
	Log.Infoln(
		"uri", uri,
		"method", method,
		"duration", duration,
	)
}

type (
	responseData struct {
		code int
		size int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.code = statusCode
}

func logResponse(code, size int) {
	Log.Infoln(
		"status code", code,
		"size", size,
	)
}

// LogMiddleware realises middleware for logging requests and responses
func LogMiddleware(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData: &responseData{
				code: 0,
				size: 0,
			},
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logRequest(uri, method, duration)
		logResponse(lw.responseData.code, lw.responseData.size)
	}

	return http.HandlerFunc(logFn)
}
