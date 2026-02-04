package database

import (
	"path/filepath"
	"testing"

	"github.com/glenntam/todoken/internal/assert"
)

func TestNew(t *testing.T) {
	t.Run("Creates DB connection pool", func(t *testing.T) {
		tmpDir := t.TempDir()
		dsn := filepath.Join(tmpDir, "db_test.sqlite")

		db, err := New(dsn)
		assert.Nil(t, err)
		assert.NotNil(t, db)
		assert.NotNil(t, db.DB)
		defer db.Close()

		err = db.Ping()
		assert.Nil(t, err)

		assert.Equal(t, 25, db.Stats().MaxOpenConnections)
	})

	t.Run("Fails with invalid DSN", func(t *testing.T) {
		dsn := "/invalid/path/that/does/not/exist/test.db"

		db, err := New(dsn)
		assert.NotNil(t, err)
		assert.Nil(t, db)
	})
}

func TestMigrateUp(t *testing.T) {
	t.Run("Applies all up migrations", func(t *testing.T) {
		db := newTestDB(t)

		err := db.MigrateUp()
		assert.Nil(t, err)

		var version int
		err = db.Get(&version, "SELECT version FROM schema_migrations LIMIT 1")
		assert.Nil(t, err)
		assert.True(t, version > 0)
	})
}
