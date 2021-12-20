package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ufal/biblio/bibdb"
)

// Start starts a REST API server.
func Start() {
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("no DB_DSN provided")
	}
	store, err := bibdb.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Print(err)
		}
	}()
	mux := http.NewServeMux()
	mux.Handle("/import-bibtex", &importBibtexHandler{store})
	server := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		Addr:              addr,
		Handler:           mux,
	}
	log.Print("going to listen on ", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
