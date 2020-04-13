package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setupFS() (tempDir string) {
	parentDir := os.TempDir()
	var err error
	tempDir, err = ioutil.TempDir(parentDir, "*-test")
	if err != nil {
		log.Fatal(err)
	}
	writeFile(tempDir + "/test.txt", []byte("hello go"))
	writeFile(tempDir + "/test.js", []byte("/* some js */"))
	return tempDir
}

func writeFile(file string, content []byte) {
	err := ioutil.WriteFile(file, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

type test struct {
	name string
	method string
	URL string
	tests func (t *testing.T, recorder *httptest.ResponseRecorder)
}

func runTests(t *testing.T, handler http.Handler, tests []test) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := http.NewRequest(test.method, test.URL, nil)
			if err != nil {
				t.Errorf("expected no error got %v", err)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			test.tests(t, w)
		})
	}
}

func cleanTempDir(tempDir string) {
	err := os.RemoveAll(tempDir) // clean up
	if err != nil {
		log.Fatalf("expected no error got %v", err)
	}
}

func TestNormalFS(t *testing.T) {
	tempDir := setupFS()
	defer cleanTempDir(tempDir)

	fs := http.FileServer(http.Dir(tempDir))

	tests := []test{
		{
			name: "Index html redirect",
			method: "GET",
			URL: "/index.html",
			tests: func (t *testing.T, rec *httptest.ResponseRecorder) {
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
			name: "Directory returns listing (no special handling of HTTP method)",
			method: "POST",
			URL: "/",
			tests: func (t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				shouldContain := "<a href=\"test.txt\">test.txt</a>"
				body := string(rec.Body.Bytes())
				if !strings.Contains(body, shouldContain) {
					t.Fatalf("%v should contain %v", body, shouldContain)
				}
			},
		},
		{
			name: "File served (no special handling of HTTP method)",
			method: "POST",
			URL: "/test.txt",
			tests: func (t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				shouldContain := "hello go"
				body := string(rec.Body.Bytes())
				if !strings.Contains(body, shouldContain) {
					t.Fatalf("%v should contain %v", body, shouldContain)
				}
			},
		},
		{
			name: "Mime type detected",
			method: "GET",
			URL: "/test.txt",
			tests: func (t *testing.T, rec *httptest.ResponseRecorder) {
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
			name: "Mime type detected (JS)",
			method: "GET",
			URL: "/test.js",
			tests: func (t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				mimeText := "application/javascript"
				location := rec.Header().Get("Content-type")
				if location != mimeText {
					t.Fatalf("Expected redirect %v but got %v", mimeText, location)
				}
			},
		},
	}

	runTests(t, fs, tests)
}

func TestRestrictedFS(t *testing.T) {
	tempDir := setupFS()
	defer cleanTempDir(tempDir)

	fs := http.FileServer(justFilesFilesystem{http.Dir(tempDir)})

	tests := []test{
		{
			name: "Index html redirect",
			method: "GET",
			URL: "/index.html",
			tests: func (t *testing.T, rec *httptest.ResponseRecorder) {
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
			name: "No directory returns listing",
			method: "GET",
			URL: "/",
			tests: func (t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				shouldBe := "<pre>\n</pre>\n"
				body := string(rec.Body.Bytes())
				if body != shouldBe {
					t.Fatalf("Expected %v but got %v", shouldBe, body)
				}
			},
		},
		{
			name: "File serving should still work as expected",
			method: "GET",
			URL: "/test.txt",
			tests: func (t *testing.T, rec *httptest.ResponseRecorder) {
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
	}

	runTests(t, fs, tests)
}

func TestRestrictedFSWithIndexHtml(t *testing.T) {
	tempDir := setupFS()
	defer cleanTempDir(tempDir)
	writeFile(tempDir + "/index.html", []byte("<html><body>Foo bar</body></html>"))

	fs := http.FileServer(justFilesFilesystem{http.Dir(tempDir)})

	tests := []test{
		{
			name: "Index html redirect",
			method: "GET",
			URL: "/index.html",
			tests: func (t *testing.T, rec *httptest.ResponseRecorder) {
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
			name: "Should serve index.html",
			method: "GET",
			URL: "/",
			tests: func (t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Fatalf("Expected %v but got %v", http.StatusOK, rec.Code)
				}
				mimeText := "text/html; charset=utf-8"
				location := rec.Header().Get("Content-type")
				if location != mimeText {
					t.Fatalf("Expected mime type %v but got %v", mimeText, location)
				}
				shouldContain := "Foo bar"
				body := string(rec.Body.Bytes())
				if !strings.Contains(body, shouldContain) {
					t.Fatalf("%v should contain %v", body, shouldContain)
				}
			},
		},
	}

	runTests(t, fs, tests)
}
