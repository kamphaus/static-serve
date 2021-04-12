package main

import (
	"bytes"
	"github.com/felixge/httpsnoop"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func LogAccess(logAccessFlag bool, prefix string, h http.Handler) http.Handler {
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
		log.Printf("%s %d %d %s", r.RemoteAddr, httpCode, writtenBytes, prefix + r.URL.Path)
	})
}

func LogReqResponse(logReqResponse bool, postfix string, h http.Handler) http.Handler {
	if !logReqResponse {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqUuid, _ := uuid.NewRandom()
		reqId := reqUuid.String()[:8] + postfix
		var (
			respCode = http.StatusOK
			hooks = httpsnoop.Hooks{
				WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						respCode = code
						next(code)
					}
				},
			}
		)
		wrapped := httpsnoop.Wrap(w, hooks)
		reqHeaders, _ := httputil.DumpRequest(r, false)
		log.Printf("%s %s", reqId, reqHeaders)
		h.ServeHTTP(wrapped, r)
		var b bytes.Buffer
		_ = wrapped.Header().WriteSubset(&b, map[string]bool{})
		log.Printf("%s %d %s", reqId, respCode, b.String())
	})
}
