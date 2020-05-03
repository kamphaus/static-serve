package main

import (
	"net/http"
)

func HandleHealthEndpoint(serveHealthEndpoint bool, h http.Handler) http.Handler {
	if !serveHealthEndpoint {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.URL.Path == "/ready" {
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
			return
		}
		h.ServeHTTP(w, r)
	})
}
