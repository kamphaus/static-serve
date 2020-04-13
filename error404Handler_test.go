package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestError404HandlerDisabled(t *testing.T) {
	tempDir := setupFS()
	defer cleanTempDir(tempDir)
	writeFile(tempDir + "/index.html", []byte("<html><body>Foo bar</body></html>"))

	error404File := ""
	h := HandleError404(&error404File, true, http.FileServer(justFilesFilesystem{http.Dir(tempDir)}))

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
					t.Fatalf("Expected redirect %v but got %v", mimeText, location)
				}
			},
		},
		{
			name: "Error 404 disabled",
			method: "GET",
			URL: "/error404.txt",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNotFound {
					t.Fatalf("Expected %v but got %v", http.StatusNotFound, rec.Code)
				}
				shouldContain := "404 page not found"
				body := string(rec.Body.Bytes())
				if !strings.Contains(body, shouldContain) {
					t.Fatalf("%v should contain %v", body, shouldContain)
				}
				mimeText := "text/plain; charset=utf-8"
				location := rec.Header().Get("Content-type")
				if location != mimeText {
					t.Fatalf("Expected redirect %v but got %v", mimeText, location)
				}
			},
		},
	}

	runLoggerTests(t, h, tests)
}

func TestError404HandlerNoLogging(t *testing.T) {
	tempDir := setupFS()
	defer cleanTempDir(tempDir)
	writeFile(tempDir + "/index.html", []byte("<html><body>Foo bar</body></html>"))

	error404File := "/"
	h := HandleError404(&error404File, false, http.FileServer(justFilesFilesystem{http.Dir(tempDir)}))

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
				log := logs.String()
				if log != "" {
					t.Fatalf("Expected '' but got %v", log)
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
			},
		},
		{
			name: "Error 404 disabled",
			method: "GET",
			URL: "/error404.txt",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				shouldContain := "Foo bar"
				body := string(rec.Body.Bytes())
				if !strings.Contains(body, shouldContain) {
					t.Fatalf("%v should contain %v", body, shouldContain)
				}
				mimeText := "text/html; charset=utf-8"
				location := rec.Header().Get("Content-type")
				if location != mimeText {
					t.Fatalf("Expected mime type %v but got %v", mimeText, location)
				}
				logStr := logs.String()
				if logStr != "" {
					t.Fatalf("Expected '' but got %v", logStr)
				}
			},
		},
	}

	runLoggerTests(t, h, tests)
}

func TestError404HandlerWithLogging(t *testing.T) {
	tempDir := setupFS()
	defer cleanTempDir(tempDir)
	writeFile(tempDir + "/index.html", []byte("<html><body>Foo bar</body></html>"))

	error404File := "/"
	h := HandleError404(&error404File, true, http.FileServer(justFilesFilesystem{http.Dir(tempDir)}))

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
				logStr := logs.String()
				if logStr != "" {
					t.Fatalf("Expected '' but got %v", logStr)
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
				logStr := logs.String()
				if logStr != "" {
					t.Fatalf("Expected '' but got %v", logStr)
				}
			},
		},
		{
			name: "Serve error page",
			method: "GET",
			URL: "/error404.txt",
			tests: func (t *testing.T, logs *bytes.Buffer, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				shouldContain := "Foo bar"
				body := string(rec.Body.Bytes())
				if !strings.Contains(body, shouldContain) {
					t.Fatalf("%v should contain %v", body, shouldContain)
				}
				mimeText := "text/html; charset=utf-8"
				location := rec.Header().Get("Content-type")
				if location != mimeText {
					t.Fatalf("Expected mime type %v but got %v", mimeText, location)
				}
				shouldContain = "Did not find /error404.txt, serving / instead"
				logStr := logs.String()
				if !strings.Contains(logStr, shouldContain) {
					t.Fatalf("%v should contain %v", logStr, shouldContain)
				}
			},
		},
	}

	runLoggerTests(t, h, tests)
}
