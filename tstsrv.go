package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

func main() {

	s := &http.Server{
		Addr:           "0.0.0.0:8443",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// add handler to default service mux
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Printf("using Server Certificate...\n")
	err := s.ListenAndServeTLS("s_cert.pem", "s_key.pem")
	if err != nil {
		log.Printf("failed to ListenAndServerTLS: %s\n", err)
	}
}
