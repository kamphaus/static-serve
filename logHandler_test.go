package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type loggerTest struct {
	name string
	method string
	URL string
	tests func (t *testing.T, logs *bytes.Buffer, recorder *httptest.ResponseRecorder)
}

func runLoggerTests(t *testing.T, handler http.Handler, tests []loggerTest) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := http.NewRequest(test.method, test.URL, nil)
			if err != nil {
				t.Errorf("expected no error got %v", err)
			}
			w := httptest.NewRecorder()
			l := log.Writer()
			buf := &bytes.Buffer{}
			log.SetOutput(buf)
			defer log.SetOutput(l)
			handler.ServeHTTP(w, r)
			test.tests(t, buf, w)
		})
	}
}

func TestDeactivatedLog(t *testing.T) {
	tempDir := setupFS()
	defer cleanTempDir(tempDir)

	h := LogAccess(false, http.FileServer(justFilesFilesystem{http.Dir(tempDir)}))

	tests := []loggerTest{
		{
			name: "Index html redirect",
			method: "GET",
			URL: "/index.html",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusMovedPermanently {
					t.Fatalf("Expected %v but got %v", http.StatusMovedPermanently, rec.Code)
				}
				location := rec.Header().Get("Location")
				if location != "./" {
					t.Fatalf("Expected redirect %v but got %v", "./", location)
				}
				expectedLog := ""
				logStr := logs.String()
				if logStr != expectedLog {
					t.Fatalf("Expected %v but got %v", expectedLog, logStr)
				}
			},
		},
		{
			name: "File serving should still work as expected",
			method: "GET",
			URL: "/test.txt",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				shouldContain := "hello go"
				body := string(rec.Body.Bytes())
				if !strings.Contains(body, shouldContain) {
					t.Fatalf("%v should contain %v", body, shouldContain)
				}
				mimeText := "text/plain; charset=utf-8"
				location := rec.Header().Get("Content-type")
				if location != mimeText {
					t.Fatalf("Expected mime type %v but got %v", mimeText, location)
				}
				expectedLog := ""
				logStr := logs.String()
				if logStr != expectedLog {
					t.Fatalf("Expected %v but got %v", expectedLog, logStr)
				}
			},
		},
	}

	runLoggerTests(t, h, tests)
}

func TestWithLogging(t *testing.T) {
	tempDir := setupFS()
	defer cleanTempDir(tempDir)

	h := LogAccess(true, http.FileServer(justFilesFilesystem{http.Dir(tempDir)}))

	tests := []loggerTest{
		{
			name: "Index html redirect",
			method: "GET",
			URL: "/index.html",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusMovedPermanently {
					t.Fatalf("Expected %v but got %v", http.StatusMovedPermanently, rec.Code)
				}
				location := rec.Header().Get("Location")
				if location != "./" {
					t.Fatalf("Expected redirect %v but got %v", "./", location)
				}
				shouldContain := "301 0 /index.html"
				logStr := logs.String()
				if !strings.Contains(logStr, shouldContain) {
					t.Fatalf("%v should contain %v", logStr, shouldContain)
				}
			},
		},
		{
			name: "File serving should still work as expected",
			method: "GET",
			URL: "/test.txt",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				shouldContain := "hello go"
				body := string(rec.Body.Bytes())
				if !strings.Contains(body, shouldContain) {
					t.Fatalf("%v should contain %v", body, shouldContain)
				}
				mimeText := "text/plain; charset=utf-8"
				location := rec.Header().Get("Content-type")
				if location != mimeText {
					t.Fatalf("Expected mime type %v but got %v", mimeText, location)
				}
				shouldContain = "200 8 /test.txt"
				logStr := logs.String()
				if !strings.Contains(logStr, shouldContain) {
					t.Fatalf("%v should contain %v", logStr, shouldContain)
				}
			},
		},
	}

	runLoggerTests(t, h, tests)
}
