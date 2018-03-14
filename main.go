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
	"os"
	"path/filepath"
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
	flag.Parse()

	docroot, err := filepath.Abs(*directory)
	if err != nil {
		log.Fatal(err)
	}
	fs := justFilesFilesystem{http.Dir(docroot)}
	http.Handle("/", http.StripPrefix("/", http.FileServer(fs)))

	log.Printf("Serving %s on HTTP port: %s\n", docroot, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
