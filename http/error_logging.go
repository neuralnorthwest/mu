package http

import (
	ht "net/http"

	"github.com/neuralnorthwest/mu/logging"
)

// statusCodeInterceptor is an HTTP middleware that intercepts the status code.
type statusCodeInterceptor struct {
	ht.ResponseWriter
	statusCode int
	path       string
	body       []byte
}

// WriteHeader intercepts the status code.
func (i *statusCodeInterceptor) WriteHeader(statusCode int) {
	i.statusCode = statusCode
	i.ResponseWriter.WriteHeader(statusCode)
}

// Write intercepts the body.
func (i *statusCodeInterceptor) Write(body []byte) (int, error) {
	// If the status code is >= 500, then the body is an error message.
	// Capture the first 512 bytes of the body.
	if i.statusCode >= 500 {
		if len(i.body) < 512 {
			i.body = append(i.body, body...)
		} else {
			i.body = append(i.body, []byte("...")...)
		}
	}
	return i.ResponseWriter.Write(body)
}

// ErrorLoggingMiddleware is an HTTP middleware that logs errors.
func ErrorLoggingMiddleware(logger logging.Logger) Middleware {
	return func(next ht.Handler) ht.Handler {
		return ht.HandlerFunc(func(w ht.ResponseWriter, r *ht.Request) {
			i := &statusCodeInterceptor{
				ResponseWriter: w,
				path:           r.URL.Path,
			}
			next.ServeHTTP(i, r)
			if i.statusCode >= 500 {
				logger.Errorw("HTTP error", "status_code", i.statusCode, "body", string(i.body))
			}
		})
	}
}
