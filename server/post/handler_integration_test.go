// +build sql

package post

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
	shorten, err := sut.saveURL("https://test.io/")
	assert.NoError(t, err)
	// cannot validate mutated data from DB
	assert.True(t, len(shorten) > 0)
}


// 3-4 ms/op on Laptop
func BenchmarkDbSave(b *testing.B) {
	//wait till database is available and migrate DB schema to latest version
	db, err := persist.ConnectAndMigrate("root@tcp(127.0.0.1:3306)/url-shortener")
	require.NoError(b, err)
	assert.NotNil(b, db)
	sut := Handler{DB: db}
	for i := 0; i < b.N; i++ {
		shorten, err := sut.saveURL("https://test.io/")
		assert.NoError(b, err)
		// cannot validate mutated data from DB
		assert.True(b, len(shorten) > 0)
	}
}
