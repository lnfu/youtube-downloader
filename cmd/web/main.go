package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/lnfu/youtube-downloader/pkg/models/postgres"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	origin   *postgres.OriginModel
	media    *postgres.MediaModel
}

func main() {
	addr := ":4000"
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	conninfo := "user=root password=secret dbname=youtube_downloader sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	defer db.Close()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		origin:   &postgres.OriginModel{DB: db},
		media:    &postgres.MediaModel{DB: db},
	}

	// listen
	server := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", addr)
	err = server.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}

	db.Close()
}
