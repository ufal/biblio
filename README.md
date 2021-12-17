# biblio

Experimental code (this is not the main repo of Biblio).

## Running the server

To start the server run:

`SERVER_ADDR=<addr> DB_DSN=<dsn> go run ./cmd/server`

The `SERVER_ADDR` variable defaults to `:8080`.

The `DB_DSN` variable is a PostgreSQL connection string.

### POST /import-bibtex

|Data|Type|Description|
|---|---|---|
|request body|BibTeX|bibliographical items in BibTeX format|
|response body|JSON|bibliographical items in JSON format + a list of known and unknown authors|

To test the BibTeX import locally run:

`cat biblio.bib | curl --data-binary "@-" http://localhost:8080/import-bibtex`

### Database

The connection string is something like

`DB_DSN=postgres://biblio:xxxxxx@localhost:5432/biblio?sslmode=disable`

### Importing authors

To import authors from an XML file run:

`DB_DSN=<dsn> go run ./cmd/xmlimporter authors.xml`
