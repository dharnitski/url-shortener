package post

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Handler proceses
type Handler struct {
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

	shortenURL := "qwe"

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
