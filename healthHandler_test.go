package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeactivatedHealthEndpoint(t *testing.T) {
	h := HandleHealthEndpoint(false, http.NotFoundHandler())

	tests := []loggerTest{
		{
			name: "Get health",
			method: "GET",
			URL: "/health",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNotFound {
					t.Fatalf("Expected %v but got %v", http.StatusNotFound, rec.Code)
				}
			},
		},
		{
			name: "Get ready",
			method: "GET",
			URL: "/ready",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNotFound {
					t.Fatalf("Expected %v but got %v", http.StatusNotFound, rec.Code)
				}
			},
		},
	}

	runLoggerTests(t, h, tests)
}

func TestActivatedHealthEndpoint(t *testing.T) {
	h := HandleHealthEndpoint(true, http.NotFoundHandler())

	tests := []loggerTest{
		{
			name: "Get health",
			method: "GET",
			URL: "/health",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNoContent {
					t.Fatalf("Expected %v but got %v", http.StatusNoContent, rec.Code)
				}
				if rec.Body.Len() != 0 {
					t.Fatalf("Expected no body but got %v", rec.Body.String())
				}
			},
		},
		{
			name: "Get health (POST)",
			method: "POST",
			URL: "/health",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNoContent {
					t.Fatalf("Expected %v but got %v", http.StatusNoContent, rec.Code)
				}
				if rec.Body.Len() != 0 {
					t.Fatalf("Expected no body but got %v", rec.Body.String())
				}
			},
		},
		{
			name: "Get ready",
			method: "GET",
			URL: "/ready",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				if rec.Body.String() != "{}" {
					t.Fatalf("Expected {} body but got %v", rec.Body.String())
				}
			},
		},
		{
			name: "Get ready (POST)",
			method: "POST",
			URL: "/ready",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				if rec.Body.String() != "{}" {
					t.Fatalf("Expected {} body but got %v", rec.Body.String())
				}
			},
		},
		{
			name: "Get ready (OPTIONS)",
			method: "OPTIONS",
			URL: "/ready",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNoContent {
					t.Fatalf("Expected %v but got %v", http.StatusNoContent, rec.Code)
				}
				if rec.Body.Len() != 0 {
					t.Fatalf("Expected no body but got %v", rec.Body.String())
				}
			},
		},
		{
			name: "Root",
			method: "GET",
			URL: "/",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNotFound {
					t.Fatalf("Expected %v but got %v", http.StatusNotFound, rec.Code)
				}
			},
		},
		{
			name: "Health prefix1",
			method: "GET",
			URL: "/healthz",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNotFound {
					t.Fatalf("Expected %v but got %v", http.StatusNotFound, rec.Code)
				}
			},
		},
		{
			name: "Health prefix2",
			method: "GET",
			URL: "/health/",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNotFound {
					t.Fatalf("Expected %v but got %v", http.StatusNotFound, rec.Code)
				}
			},
		},
	}

	runLoggerTests(t, h, tests)
}
