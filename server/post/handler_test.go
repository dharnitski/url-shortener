package post_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dharnitski/url-shortener/post"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostHandler(t *testing.T) {
	sut := post.Handler{}

	r, err := http.NewRequest(http.MethodPost, "https://example.com/", bytes.NewBufferString(`{"url": "https://github.com/"}`))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	assert.Equal(t, `{"shortenUrl":"qwe"}`, string(body))
}

func TestPostHandlerInvalidJson(t *testing.T) {
	sut := post.Handler{}

	r, err := http.NewRequest(http.MethodPost, "https://example.com/", bytes.NewBufferString(`invalid`))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut.ServeHTTP(w, r)
	assert.Equal(t, 400, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	assert.Equal(t, "invalid character 'i' looking for beginning of value\n", string(body))
}

func TestPostHandlerEmptyURL(t *testing.T) {
	sut := post.Handler{}

	r, err := http.NewRequest(http.MethodPost, "https://example.com/", bytes.NewBufferString(`{}`))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut.ServeHTTP(w, r)
	assert.Equal(t, 400, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	assert.Equal(t, "empty URL\n", string(body))
}
