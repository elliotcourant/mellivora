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

func TestGetModelInfo(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		type DataNode struct {
			DataNodeId uint64 `m:"pk,serial"`
			Address    string `m:"unique:uq_address_port"`
			Port       int32  `m:"unique:uq_address_port"`
			User       string
			Password   string
			Healthy    bool
		}
		info := getModelInfo(DataNode{})
		assert.Equal(t, "DataNode", info.Name())
	})
}
