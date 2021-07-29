package main

import (
	"crypto/tls"
	"log"
)

func loadTlsConfig(verbose bool, tlsCertFlag string, tlsKeyFlag string) *tls.Config {
	if tlsCertFlag == "" && tlsKeyFlag == "" {
		return nil
	}
	if tlsCertFlag == "" {
		log.Fatal("Path to TLS key file ist set. Path to certificate file is required, but missing.")
	}
	if tlsKeyFlag == "" {
		log.Fatal("Path to TLS certificate file ist set. Path to certificate key file is required, but missing.")
	}
	cert, err := tls.LoadX509KeyPair(tlsCertFlag, tlsKeyFlag)
	if err != nil {
		log.Fatal(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}
