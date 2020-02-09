// +build sql

package get

import (
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dharnitski/url-shortener/persist"
)

func TestDB(t *testing.T) {
	//wait till database is available and migrate DB schema to latest version
	db, err := persist.ConnectAndMigrate("root@tcp(127.0.0.1:3306)/url-shortener")
	require.NoError(t, err)
	assert.NotNil(t, db)
	sut := Handler{DB: db}
	_, err = sut.getURL(42)
	assert.NoError(t, err)
	// validation for data is not necessary. It is tested by unit test
	// error checked to make sure that SQL matches DB schema
}
