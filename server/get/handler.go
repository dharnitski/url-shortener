// Package get implements handler that returns redirects to full url or 404 for shorten path
package get

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dharnitski/url-shortener/shortener"
	"github.com/gorilla/mux"
)

// Handler proceses
type Handler struct {
	DB *sql.DB
	// App hosting domain with protocol, port, and slash - http://localhost:8080/
	UIDomain string
}

// ServeHTTP serves requests
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shorten := vars["shorten"]

	// no need to validate for empty param - checked by router

	// decode in64 from shorten URL
	id, err := shortener.ShortToInt(shorten)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path, err := h.getURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if path == "" {
		http.Error(w, fmt.Sprintf("No url for %q", shorten), http.StatusNotFound)
		return
	}

	// Return temporary redirect to enable landing URL chanage in the future
	http.Redirect(w, r, fmt.Sprintf("%s%s", h.UIDomain, path), http.StatusFound)
}

// getURL returns full url from DB
// empty string return if there is no record with given id
func (h Handler) getURL(id int64) (string, error) {
	var url string
	row := h.DB.QueryRow("select url from links where id = ?", id)
	switch err := row.Scan(&url); err {
	case sql.ErrNoRows:
		// no matching record
		return "", nil
	case nil:
		// success
		return url, nil
	default:
		// error
		return "", err
	}
}
