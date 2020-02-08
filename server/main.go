package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"

	"github.com/dharnitski/url-shortener/post"
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func connectAndMigrate(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return db, err
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return db, err
	}

	migration, err := migrate.NewWithDatabaseInstance("file://./migrations", "mysql", driver)
	if err != nil {
		return db, err
	}

	err = migration.Up()
	if err != nil {
		// should not fail if Database Schema is up to date or it is currently updated (and locked) by other instance
		if err == migrate.ErrNoChange || err == migrate.ErrLocked {
			fmt.Print("DB Migration: " + err.Error())
			return db, nil
		}
		return db, err
	}

	return db, nil
}

func main() {
	connectionString, exists := os.LookupEnv("DB_CONNECTION_STRING")
	if !exists {
		log.Fatal("requires env variable DB_CONNECTION_STRING")
	}

	// db implements connection pool and it is safe to use in multiple goroutines
	db, err := connectAndMigrate(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.Handle("/api/post", post.Handler{DB: db}).Methods("POST")

	spa := spaHandler{staticPath: "ui", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
