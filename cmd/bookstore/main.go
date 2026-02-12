package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/sergekukharev/agent-test-writer-validator/internal/api"
	"github.com/sergekukharev/agent-test-writer-validator/internal/storage"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	repo := storage.NewBookRepository()
	handler := api.NewHandler(repo)
	mux := handler.Routes()

	var h http.Handler = mux
	h = api.LoggingMiddleware(h)
	h = api.RecoveryMiddleware(h)

	log.Printf("bookstore listening on %s", *addr)
	if err := http.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
