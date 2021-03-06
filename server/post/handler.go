package post

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dharnitski/url-shortener/shortener"
)

// Handler proceses
type Handler struct {
	DB *sql.DB
	// App hosting domain with protocol, port, and slash - http://localhost:8080/
	UIDomain string
}

// Request contains URL to shorten
type Request struct {
	URL string `json:"url"`
}

// Response contains ShortenURL
type Response struct {
	ShortenURL string `json:"shortenUrl"`
}

// ServeHTTP serves requests
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req Request

	// encode JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate that URL is not empty
	if req.URL == "" {
		http.Error(w, "empty URL", http.StatusBadRequest)
		return
	}

	// validate if URL is valid
	parsed, err := url.Parse(req.URL)
	if err != nil || parsed.Host == "" {
		http.Error(w, "not valid URL", http.StatusBadRequest)
		return
	}

	// validate URL scheme
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		http.Error(w, "only web links supported", http.StatusBadRequest)
		return
	}

	// infinite redirect loop prevention
	// case insensitive strings comparison to match GOOGLE.COM and google.com
	if strings.HasPrefix(strings.ToLower(req.URL), strings.ToLower(h.UIDomain)) {
		http.Error(w, "links to this site are not supported to prevent infinite redirects", http.StatusBadRequest)
		return
	}

	shortenURL, err := h.saveURL(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		ShortenURL: shortenURL,
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h Handler) saveURL(link string) (string, error) {
	result, err := h.DB.Exec("INSERT INTO links (url) VALUES(?)", link)
	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	shortened := shortener.IntToShort(id)
	return fmt.Sprintf("%s%s", h.UIDomain, shortened), nil
}
