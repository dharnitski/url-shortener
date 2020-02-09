package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"

	"github.com/dharnitski/url-shortener/get"
	"github.com/dharnitski/url-shortener/persist"
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

func main() {
	connectionString, exists := os.LookupEnv("DB_CONNECTION_STRING")
	if !exists {
		log.Fatal("requires env variable DB_CONNECTION_STRING")
	}
	// App hosting domain with protocol, port, and slash - http://localhost:8080/
	uiDomain, exists := os.LookupEnv("UI_DOMAIN")
	if !exists {
		log.Fatal("requires env variable UI_DOMAIN")
	}

	// db implements connection pool and it is safe to use in multiple goroutines
	db, err := persist.ConnectAndMigrate(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	// Rest API register new URL
	router.Handle("/api/post", post.Handler{DB: db, UIDomain: uiDomain}).Methods("POST")
	// HTTP handler to handle short urls
	router.Handle("/{shorten:[0-9a-zA-Z]+}", get.Handler{DB: db}).Methods("GET")

	// host SPA application
	// important!!! SPA and short url handler shares the same subroute
	// implementation works based on assumption that no spa resources match short url RegEx
	// that is solid for current version of client side compiler
	// it compiles all the files to assets with file extention like some_script.js or favicon.ico
	// assets has to be moved to different folder if for any reason they match short url RegEx
	// Issue will be exposed as Redirect or NotFound responses to asset routes 
	spa := spaHandler{staticPath: "ui", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
