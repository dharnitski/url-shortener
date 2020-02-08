package post_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dharnitski/url-shortener/post"
)

func TestPostHandler(t *testing.T) {
	t.Parallel()
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
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			sut := post.Handler{}

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
