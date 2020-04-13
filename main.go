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
	"os"
	"os/signal"
	"path/filepath"
	"sync"
)

type arrayFlags []string
func (i *arrayFlags) String() string {
	return "my string representation"
}
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
var ports arrayFlags
var directories arrayFlags
var error404s arrayFlags

func main() {
	log.SetFlags(0)
	flag.Var(&ports, "p", "ports to serve on (default: 8100)")
	flag.Var(&directories, "d", "ports to serve on (default: ./)")
	flag.Var(&error404s, "e", "the files to serve in case of error 404 (- to disable error404 handler)")
	logAccessFlag := flag.Bool("l", false, "log access requests")
	error404VerboseFlag := flag.Bool("v", false, "log when handling error 404")
	flag.Parse()

	if len(ports) != len(directories) || len(ports) != len(error404s) {
		flag.Usage()
		log.Printf("Ports, directories and error404 flags need to be specified the same amount of times.")
		os.Exit(1)
		return
	}
	if len(ports) == 0 {
		ports = append(ports, "8100")
		directories = append(directories, ".")
		error404s = append(error404s, "-")
	}

	var servers []*http.Server
	var done sync.WaitGroup
	done.Add(len(ports))
	for i := range ports {
		port := ports[i]
		directory := directories[i]
		error404File := error404s[i]
		if error404File == "-" {
			error404File = ""
		}
		servers = append(servers, serve(&done, port, directory, error404File, len(ports), *logAccessFlag, *error404VerboseFlag))
	}

	// run until we get a signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	log.Printf("Shutting down...")
	for _, server := range servers {
		err := server.Close()
		if err != nil {
			log.Printf("Encountered error: %v", err)
		}
	}
	done.Wait()
	log.Printf("Shutdown complete")
}

func serve(wg *sync.WaitGroup, port string, directory string, error404File string, numPorts int, logAccess bool, error404Verbose bool) *http.Server {
	docroot, err := filepath.Abs(directory)
	if err != nil {
		log.Fatal(err)
	}
	fs := justFilesFilesystem{http.Dir(docroot)}

	withError404 := ""
	if error404File != "" {
		withError404 = fmt.Sprintf(" with %s as error 404 file", error404File)
	}
	log.Printf("Serving %s on HTTP port: %s%s\n", docroot, port, withError404)
	listenAddr := ":" + port
	logPrefix := listenAddr
	if numPorts == 1 {
		logPrefix = ""
	}
	server := &http.Server{Addr: listenAddr, Handler: LogAccess(logAccess, logPrefix, HandleError404(&error404File, error404Verbose, http.StripPrefix("/", http.FileServer(fs))))}
	go func() {
		defer func() { wg.Done() }()
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("Encountered error: %v", err)
		}
	}()
	return server
}
