// +build sql

package persist

import (
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	//wait till database is available and migrate DB schema to latest version
	db, err := ConnectAndMigrate("root@tcp(127.0.0.1:3306)/url-shortener")
	require.NoError(t, err)
	assert.NotNil(t, db)
}
