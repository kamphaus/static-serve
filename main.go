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
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/ranveerkunal/memfs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
)

type arrayFlags []string
func (i *arrayFlags) String() string {
	return strings.Join(*i,",")
}
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type FSType string
const (
	DiskFS FSType = "diskfs"
	INMem = "inmem"
	INMemWithoutWatch = "inmem-nowatch"
)
var fsTypes = []FSType{DiskFS, INMem, INMemWithoutWatch}
var invalidFSType = errors.New("Invalid filesystem type")
func (i *FSType) String() string {
	return string(*i)
}
func (i *FSType) Set(value string) error {
	input := FSType(strings.ToLower(value))
	for _, val := range fsTypes {
		if val == input {
			*i = val
			return nil
		}
	}
	return invalidFSType
}

var ports arrayFlags
var directories arrayFlags
var error404s arrayFlags
var fsType = DiskFS

func main() {
	log.SetFlags(0)
	flag.Var(&ports, "p", "ports to serve on (default: 8100)")
	flag.Var(&directories, "d", "the directories of static files to host (default: ./)")
	flag.Var(&error404s, "e", "the files to serve in case of error 404 (- to disable error404 handler)")
	flag.Var(&fsType, "fs-type", "Which filesystem type to use. Options:\n" +
		"* "+string(DiskFS)+"        Load files directly from disk (Kernel takes care of caching)\n" +
		"* "+string(INMem)+"         Eagerly loads files from directories into memory and serves them from memory\n" +
		"* "+string(INMemWithoutWatch)+" Same as "+string(INMem)+", but doesn't watch for changes (ideal for docker containers)\n")
	logAccessFlag := flag.Bool("l", false, "log access requests")
	error404VerboseFlag := flag.Bool("v", false, "log when handling error 404")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		log.Printf("Ports, directories and error404 flags can be specified multiple times, but need to be specified the same amount of times.")
	}
	flag.Parse()

	if len(ports) != len(directories) || len(ports) != len(error404s) {
		flag.Usage()
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
		servers = append(servers, serve(&done, port, directory, error404File, len(ports), fsType, *logAccessFlag, *error404VerboseFlag))
	}

	// run until we get a signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	log.Printf("Shutting down...")
	for _, server := range servers {
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Printf("HTTP server Shutdown error: %v", err)
		}
	}
	done.Wait()
	log.Printf("Shutdown complete")
}

func serve(wg *sync.WaitGroup, port string, directory string, error404File string, numPorts int, fsType FSType, logAccess bool, error404Verbose bool) *http.Server {
	docroot, err := filepath.Abs(directory)
	if err != nil {
		log.Fatal(err)
	}
	var fs http.FileSystem
	if fsType == INMem {
		fs, err = memfs.NewWithWatch(docroot, true)
	} else if fsType == INMemWithoutWatch {
		fs, err = memfs.NewWithWatch(docroot, false)
	} else {
		fs = http.Dir(docroot)
	}
	fs = justFilesFilesystem{fs}

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
