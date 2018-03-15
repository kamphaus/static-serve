/*
Serve is a very simple static file server in go
Usage:
	-p="8100": port to serve on
	-d=".":    the directory of static files to host

Serve does not show directory listings, it only serves files.
*/
package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/felixge/httpsnoop"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func main() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	error404File := flag.String("e", "", "the file to serve in case of error 404")
	flag.Parse()

	docroot, err := filepath.Abs(*directory)
	if err != nil {
		log.Fatal(err)
	}
	fs := justFilesFilesystem{http.Dir(docroot)}
	http.Handle("/", HandleError404(error404File, http.StripPrefix("/", http.FileServer(fs))))

	withError404 := ""
	if error404File != nil && *error404File != "" {
		withError404 = fmt.Sprintf(" with %s as error 404 file", *error404File)
	}
	log.Printf("Serving %s on HTTP port: %s%s\n", docroot, *port, withError404)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

var ignoreError404 = errors.New("ignored file")

func HandleError404(error404File *string, h http.Handler) http.Handler {
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
			log.Printf("Did not find %s, serving %s instead", r.URL.Path, *error404File)
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
