package mellivora

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetBaseTypeOf(t *testing.T) {
	type Simple struct {
		Field int
	}

	t.Run("simple", func(t *testing.T) {
		typ := getBaseTypeOf(Simple{})
		assert.Equal(t, "Simple", typ.Name())
		assert.Equal(t, reflect.Struct, typ.Kind())
	})

	t.Run("slice", func(t *testing.T) {
		typ := getBaseTypeOf([]Simple{})
		assert.Equal(t, "Simple", typ.Name())
		assert.Equal(t, reflect.Struct, typ.Kind())
	})

	t.Run("invalid", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = getBaseTypeOf(1)
		})
	})
}
