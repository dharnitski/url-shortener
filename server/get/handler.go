package get

import (
	"fmt"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dharnitski/url-shortener/shortener"
	"github.com/gorilla/mux"
)

// Handler proceses
type Handler struct {
	DB *sql.DB
}

// Response contains ShortenURL
type Response struct {
	URL string `json:"url"`
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

	url, err := h.getURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if url == "" {
		http.Error(w, fmt.Sprintf("No url for %q", shorten), http.StatusNotFound)
		return
	}

	response := Response{
		URL: url,
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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
