package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Log represents global var for logging.
// By default Log is no-op logger.
var Log *zap.SugaredLogger = zap.NewNop().Sugar()

// Initialize creates logger Log.
func Initialize() error {
	zl, err := zap.NewProduction()
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

func logResponse(code, size int) {
	Log.Infoln(
		"status code", code,
		"size", size,
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

// Write redefines Write method of http.ResponseWriter.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader redefines WriteHeader method of http.ResponseWriter.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.code = statusCode
}

// LogMiddleware realises middleware for logging requests and responses.
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
