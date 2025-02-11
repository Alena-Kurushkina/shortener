// Package compress realises middleware for responces compression and requests decompression.
package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressWriter defines object for compressing output responces.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter construct compressWriter.
func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header redefines func Header of http.ResponseWriter.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// WriteHeader redefines func WriteHeader of http.ResponseWriter.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Write redefines func Write of http.ResponseWriter.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// Close redefines func Close of http.ResponseWriter.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader defines object for decompressing income requests.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader construct compressReader.
func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read redefine Read method of io.ReadCloser.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close redefine Close method of io.ReadCloser.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// GzipMiddleware realises middleware for requests and responses compression in gzip format.
// It compresses such content types as "application/json", "text/html", "application/x-gzip".
// It avoid compression if Accept-Encoding or Content-Encoding headers which don`t contain "gzip".
func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		supportsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		validContentTypes := []string{"application/json", "text/html", "application/x-gzip"}
		isValid := false
		contentType := r.Header.Get("Content-Type")
		for _, cntType := range validContentTypes {
			if strings.Contains(contentType, cntType) {
				isValid = true
			}
		}

		if supportsGzip && isValid {
			cw := NewCompressWriter(w)
			ow = cw
			defer func(){
				tErr:=cw.Close()
				if tErr != nil {
					http.Error(w, tErr.Error(), http.StatusInternalServerError)
					return
				}
			}()
		}

		sendsGzip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if sendsGzip {
			cr, err := NewCompressReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func(){
				tErr:=cr.Close()
				if tErr != nil {
					http.Error(w, tErr.Error(), http.StatusInternalServerError)
					return
				}
			}()
		}

		h.ServeHTTP(ow, r)
	})
}
