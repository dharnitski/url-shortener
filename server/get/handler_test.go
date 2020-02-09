package get_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dharnitski/url-shortener/get"
)

func TestPostHandler(t *testing.T) {
	t.Parallel()
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	rows := sqlmock.NewRows([]string{"url"}).AddRow("G")
	sqlMock.ExpectQuery(`select url from links where id = ?`).
		WithArgs(42).
		WillReturnRows(rows)
	sut := get.Handler{DB: db, UIDomain: "http://localhost:8080/"}

	router := mux.NewRouter()
	router.Handle("/{shorten}", sut)

	r, err := http.NewRequest(http.MethodGet, "https://example.com/G", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	assert.Equal(t, 302, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	assert.Equal(t, "<a href=\"http://localhost:8080/G\">Found</a>.\n\n", string(body))
}

func TestPostHandler404(t *testing.T) {
	t.Parallel()
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	rows := sqlmock.NewRows([]string{"url"}) // .AddRow("G") - db return no records
	sqlMock.ExpectQuery(`select url from links where id = ?`).
		WithArgs(42).
		WillReturnRows(rows)
	sut := get.Handler{DB: db}

	router := mux.NewRouter()
	router.Handle("/{shorten}", sut)

	r, err := http.NewRequest(http.MethodGet, "https://example.com/G", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	assert.Equal(t, 404, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	assert.Equal(t, "No url for \"G\"\n", string(body))
}

func TestPostHandlerDbError(t *testing.T) {
	t.Parallel()
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectQuery(`select url from links where id = ?`).
		WithArgs(42).
		WillReturnError(errors.New("DB Error"))
	sut := get.Handler{DB: db}

	router := mux.NewRouter()
	router.Handle("/{shorten}", sut)

	r, err := http.NewRequest(http.MethodGet, "https://example.com/G", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	assert.Equal(t, 500, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	assert.Equal(t, "DB Error\n", string(body))
}

func TestPostHandler400(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{
			"___",
			"invalid character '_' at position 0 (0)\n",
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			// should return before DB used
			sut := get.Handler{DB: nil}

			router := mux.NewRouter()
			router.Handle("/api/get/{shorten}", sut)

			r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://example.com/api/get/%s", test.input), nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)
			assert.Equal(t, 400, w.Code)
			body, err := ioutil.ReadAll(w.Body)
			require.NoError(t, err)
			assert.Equal(t, test.expected, string(body))
		})
	}
}
