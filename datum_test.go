package mellivora

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDatumBuilderBase_Keys(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		type DataNode struct {
			DataNodeId uint64 `m:"pk,serial"`
			Address    string `m:"unique:uq_address_port"`
			Port       int32  `m:"unique:uq_address_port"`
			User       string
			Password   string
			Healthy    bool
		}

		dataNode := DataNode{
			DataNodeId: 1234,
			Address:    "127.0.0.1",
			Port:       5432,
			User:       "POSTGRES",
			Password:   "password",
			Healthy:    true,
		}

		info := getModelInfo(dataNode)

		builder := newDatumBuilder(info, reflect.ValueOf(dataNode))
		datums := builder.Keys()
		assert.NotEmpty(t, datums)

		verify := builder.Verify()
		assert.NotEmpty(t, verify)
	})
}
