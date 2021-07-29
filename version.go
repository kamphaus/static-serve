package main

import (
	"log"
	"runtime"
)

var (
	// Version holds the current version of static-serve.
	version = "dev"
	// Commit holds the commit used to build static-serve.
	commit = "abcde123"
	// BuildDate holds the build date of static-serve.
	date = "I don't remember exactly"
	// BuiltBy holds who built static-serve.
	builtBy = "someone"
)

func printVersion() {
	log.Printf("  Version:    " + version +
		"\n  GoVersion:  " + runtime.Version() +
		"\n  Commit:     " + commit +
		"\n  Build time: " + date +
		"\n  Built by:   " + builtBy +
		"\n  OS:         " + runtime.GOOS +
		"\n  Arch:       " + runtime.GOARCH + "\n")
}
