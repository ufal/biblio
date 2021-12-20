package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ufal/biblio"
	"github.com/ufal/biblio/bibdb"
	"github.com/ufal/biblio/xmldb"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "no input files")
		os.Exit(1)
	}
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
	files := make(map[string]*xmldb.File)
	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}()
		file, err := xmldb.LoadFile(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		files[filepath.Base(flag.Arg(i))] = file
	}
	for name, file := range files {
		switch name {
		case "authors.xml":
			for _, record := range file.Records {
				author, err := biblio.NewAuthorFromRecord(record)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				a, err := store.GetAuthorByID(author.ID)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if a == nil {
					fmt.Println("inserting", author.FullName())
					if err := store.InsertAuthor(author); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
				} else {
					fmt.Println("ignoring", author.FullName())
				}
			}
		}
	}
}
