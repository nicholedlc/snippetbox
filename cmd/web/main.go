package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.org/pkg/models"
)

func main() {
	addr := flag.String("addr", ":4040", "HTTP network address")
	dsn := flag.String("dsn", "ndlc:secret@/snippetbox?parseTime=true", "MySQL DSN")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")

	flag.Parse()

	dbConnection := connect(*dsn)

	defer dbConnection.Close()

	app := &App{
		Database:  &models.Database{dbConnection},
		HTMLDir:   *htmlDir,
		StaticDir: *staticDir,
	}

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, app.Routes())
	log.Fatal(err)
}

// The connect() function wraps sql.Open() and returns a sql.DB connection pool for a given DSN
func connect(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
