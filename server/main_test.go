package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServeSpa404(t *testing.T) {
	sut := spaHandler{
		staticPath: "not_existing",
		indexPath:  "index.html",
	}

	r, err := http.NewRequest(http.MethodGet, "https://example.com/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut.ServeHTTP(w, r)
	assert.Equal(t, 404, w.Code)
}

func TestServeSpa(t *testing.T) {
	sut := spaHandler{
		staticPath: "testdata",
		indexPath:  "index.html",
	}

	r, err := http.NewRequest(http.MethodGet, "https://example.com/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	assert.Equal(t, `<html></html>`, string(body))
}
