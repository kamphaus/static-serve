package main

import (
	"github.com/felixge/httpsnoop"
	"io"
	"log"
	"net/http"
)

func LogAccess(logAccessFlag bool, h http.Handler) http.Handler {
	if !logAccessFlag {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			httpCode = http.StatusOK
			writtenBytes int64 = 0
			hooks = httpsnoop.Hooks{
				WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						httpCode = code
						next(code)
					}
				},
				Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return func(p []byte) (int, error) {
						n, err := next(p)
						writtenBytes += int64(n)
						return n, err
					}
				},
				ReadFrom: func(fromFunc httpsnoop.ReadFromFunc) httpsnoop.ReadFromFunc {
					return func(src io.Reader) (int64, error) {
						n, err := fromFunc(src)
						writtenBytes += n
						return n, err
					}
				},
			}
		)
		wrapped := httpsnoop.Wrap(w, hooks)
		h.ServeHTTP(wrapped, r)
		log.Printf("%s %d %d %s", r.RemoteAddr, httpCode, writtenBytes, r.URL.Path)
	})
}
