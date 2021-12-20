package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ufal/biblio/bibdb"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "no DB_DSN provided")
		os.Exit(1)
	}
	store, err := bibdb.NewStore(dsn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Print(err)
		}
	}()
	fmt.Println("DB opened")
}
