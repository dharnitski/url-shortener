package post_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dharnitski/url-shortener/post"
)

func TestPostHandler(t *testing.T) {
	t.Parallel()
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectExec(`INSERT INTO links \(url\) VALUES\(\?\)`).
		WithArgs("https://github.com/").
		WillReturnResult(sqlmock.NewResult(42, 1))
	sut := post.Handler{DB: db, UIDomain: "http://localhost:8080/"}

	r, err := http.NewRequest(http.MethodPost, "https://example.com/", bytes.NewBufferString(`{"url": "https://github.com/"}`))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	sut.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	assert.Equal(t, `{"shortenUrl":"http://localhost:8080/G"}`, string(body))
}

func TestPostHandlerInvalidInput(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected string
	}{
		{
			"invalid",
			"invalid character 'i' looking for beginning of value\n",
		},
		{
			"{}",
			"empty URL\n",
		},
		{
			`{"url": "not url"}`,
			"not valid URL\n",
		},
		{
			`{"url": "😊"}`,
			"not valid URL\n",
		},
		{
			`{"url": "ftp://github.com/"}`,
			"only web links supported\n",
		},
		{
			`{"url": "http://localhost:8080/some"}`,
			"links to this site are not supported to prevent infinite redirects\n",
		},
		{
			`{"url": "http://locAlhost:8080/some"}`,
			"links to this site are not supported to prevent infinite redirects\n",
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			sut := post.Handler{UIDomain: "http://localhost:8080/"}

			r, err := http.NewRequest(http.MethodPost, "https://example.com/", bytes.NewBufferString(test.input))
			require.NoError(t, err)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)
			assert.Equal(t, 400, w.Code)
			body, err := ioutil.ReadAll(w.Body)
			require.NoError(t, err)
			assert.Equal(t, test.expected, string(body))
		})
	}
}
