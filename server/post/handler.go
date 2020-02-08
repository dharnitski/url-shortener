package post

import (
	"encoding/json"
	"net/http"
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

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "empty URL", http.StatusBadRequest)
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
