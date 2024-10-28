package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w: w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Write(p []byte) (int, error){
	return c.zw.Write(p)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// http.Request's Body has type io.ReadCloser
type compressReader struct {
	r io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err:=gzip.NewReader(r)
	if err!=nil{
		return nil, err
	}

	return &compressReader{
		r: r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error){
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err:=c.r.Close(); err!=nil{
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		ow:=w

		supportsGzip:=strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		if supportsGzip {
			cw:=NewCompressWriter(w)
			ow=cw
			defer cw.Close()
		}

		sendsGzip:=strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if sendsGzip {
			cr, err:=NewCompressReader(r.Body)
			if err!=nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body=cr
			defer cr.Close()
		}

		h.ServeHTTP(ow,r)
	})
}