package main

import (
	"log"
	"runtime"
)

var (
	// Version holds the current version of static-serve.
	Version = "dev"
	// BuildDate holds the build date of static-serve.
	BuildDate = "I don't remember exactly"
)

func printVersion() {
	log.Printf("  Version:   " + Version +
		"\n  GoVersion: " + runtime.Version() +
		"\n  BuildTime: " + BuildDate +
		"\n  OS:        " + runtime.GOOS +
		"\n  Arch:      " + runtime.GOARCH + "\n")
}
