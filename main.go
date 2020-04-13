/*
Static-serve is a very simple static file server in go
Usage:
	-p="8100": port to serve on
	-d=".":    the directory of static files to host
	-e="":     a file to serve in case of error 404
	-l:        log access requests

Static-serve does not show directory listings, it only serves files.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	error404File := flag.String("e", "", "the file to serve in case of error 404")
	logAccessFlag := flag.Bool("l", false, "log access requests")
	flag.Parse()

	docroot, err := filepath.Abs(*directory)
	if err != nil {
		log.Fatal(err)
	}
	fs := justFilesFilesystem{http.Dir(docroot)}
	http.Handle("/", LogAccess(*logAccessFlag, HandleError404(error404File, *logAccessFlag, http.StripPrefix("/", http.FileServer(fs)))))

	withError404 := ""
	if error404File != nil && *error404File != "" {
		withError404 = fmt.Sprintf(" with %s as error 404 file", *error404File)
	}
	log.Printf("Serving %s on HTTP port: %s%s\n", docroot, *port, withError404)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
