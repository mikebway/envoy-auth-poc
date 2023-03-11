package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	// requestCount tracks the number of requests received
	requestCount int
)

// main is the command line entry point to the application.
func main() {

	// Route the annoying requests for favicon.ico
	http.HandleFunc("/favicon.ico", handleFavicon)

	// Route request for the root URL path and everything else under it
	http.HandleFunc("/", handleRoot)

	// Start the server
	const serverAddress = ":9090"
	log.Printf("starting server at %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

// handleRoot is the primary HTTP request handler for the root path and everything under it for this very crude website.
func handleRoot(w http.ResponseWriter, r *http.Request) {

	// Bump the count of requests seen
	requestCount++

	// Log the URL that we have been asked for etc
	log.Printf("handling request = %s, request count = %d\n", r.URL.Path, requestCount)

	// Dump some text onto the response
	_, _ = fmt.Fprintf(w, "Path:\t\t%q\n", r.URL.Path)
	_, _ = fmt.Fprintf(w, "Count:\t\t%d\n", requestCount)
}

// handleFavicon is the HTTP request handler for "/favicon.ico" requests: 404 NOT FOUND it!
func handleFavicon(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request = %s\n", r.URL.Path)
	http.Error(w, "Not Found", http.StatusNotFound)
}

// requestCounter tracks the number of requests, logging the count has it increases.
func requestCounter() {
	log.Printf("handling root request. Count=%d", requestCount)
}
