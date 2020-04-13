package main

import (
	"errors"
	"github.com/felixge/httpsnoop"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var ignoreError404 = errors.New("ignored file")

func HandleError404(error404File *string, shouldLog bool, h http.Handler) http.Handler {
	if error404File == nil || *error404File == "" {
		return h
	}
	if !strings.HasPrefix(*error404File, "/") {
		*error404File = "/" + *error404File
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			isError404 = false
			hooks = httpsnoop.Hooks{
				WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						if code == http.StatusNotFound {
							isError404 = true
						} else {
							next(code)
						}
					}
				},
				Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return func(p []byte) (int, error) {
						if !isError404 {
							n, err := next(p)
							return n, err
						}
						return 0, ignoreError404
					}
				},
			}
			originalHeader http.Header
		)
		wrapped := httpsnoop.Wrap(w, hooks)
		CopyHeaders(originalHeader, w.Header())
		h.ServeHTTP(wrapped, r)
		if isError404 {
			if shouldLog {
				log.Printf("Did not find %s, serving %s instead", r.URL.Path, *error404File)
			}
			SetHeaders(w.Header(), originalHeader)
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = *error404File
			r2.RequestURI = r2.URL.RequestURI()
			h.ServeHTTP(w, r2)
		}
	})
}

// CopyHeaders copies http headers from source to destination, it
// does not override, but adds multiple headers
func CopyHeaders(dst http.Header, src http.Header) {
	for k, vv := range src {
		dst[k] = append(dst[k], vv...)
	}
}

func SetHeaders(dst http.Header, src http.Header) {
	for h := range dst {
		dst.Del(h)
	}
	CopyHeaders(dst, src)
}
