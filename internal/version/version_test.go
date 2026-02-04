package version

import (
	"testing"

	"github.com/glenntam/todoken/internal/assert"
)

func TestGet(t *testing.T) {
	t.Run("Returns a non-empty string", func(t *testing.T) {
		version := Get()
		assert.True(t, version != "")
	})
}
