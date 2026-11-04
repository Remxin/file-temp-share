package main

import (
	"fmt"
	"os"
)

// Stub entrypoint for the HTTP API.
// Wire Gin, config, DB, handlers here in later iterations.
func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	fmt.Printf("file-temp-share api stub — listening placeholder on %s (not started yet)\n", addr)
}
