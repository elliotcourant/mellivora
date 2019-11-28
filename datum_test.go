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

		builder := newDatumBuilder(info, reflect.ValueOf(dataNode), true)
		datums, err := builder.Keys()
		assert.NoError(t, err)
		assert.NotEmpty(t, datums)

		verify, err := builder.Verify()
		assert.NoError(t, err)
		assert.NotEmpty(t, verify)
	})

	t.Run("array", func(t *testing.T) {
		type DataNode struct {
			DataNodeId uint64 `m:"pk,serial"`
			Address    string `m:"unique:uq_address_port"`
			Port       int32  `m:"unique:uq_address_port"`
			User       string
			Password   string
			Healthy    bool
		}

		dataNode := []*DataNode{
			{
				DataNodeId: 1234,
				Address:    "127.0.0.1",
				Port:       5432,
				User:       "POSTGRES",
				Password:   "password",
				Healthy:    true,
			},
			{
				DataNodeId: 1235,
				Address:    "127.0.0.1",
				Port:       5433,
				User:       "POSTGRES",
				Password:   "password",
				Healthy:    true,
			},
		}

		info := getModelInfo(dataNode)

		builder := newDatumBuilder(info, reflect.ValueOf(dataNode), true)
		datums, err := builder.Keys()
		assert.NoError(t, err)
		assert.NotEmpty(t, datums)

		verify, err := builder.Verify()
		assert.NoError(t, err)
		assert.NotEmpty(t, verify)
	})

	t.Run("unique", func(t *testing.T) {
		type DataNode struct {
			DataNodeId uint64 `m:"pk,serial"`
			Address    string `m:"unique:uq_address_port"`
			Port       int32  `m:"unique:uq_address_port"`
			User       string
			Password   string
			Healthy    bool
		}

		dataNode := []*DataNode{
			{
				DataNodeId: 1234,
				Address:    "127.0.0.1",
				Port:       5432,
				User:       "POSTGRES",
				Password:   "password",
				Healthy:    true,
			},
			{
				DataNodeId: 1235,
				Address:    "127.0.0.1",
				Port:       5432,
				User:       "POSTGRES",
				Password:   "password",
				Healthy:    true,
			},
		}

		info := getModelInfo(dataNode)

		builder := newDatumBuilder(info, reflect.ValueOf(dataNode), true)
		_, err := builder.Keys()
		assert.Error(t, err, "expected error due to unique violation")
	})
}

func TestDatumReaderBase_Read(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		type Item struct {
			ItemId uint64 `m:"pk"`
			Name   string
		}

		items := []Item{
			{
				ItemId: 1,
				Name:   "Item One",
			},
			{
				ItemId: 2,
				Name:   "Item Two",
			},
			{
				ItemId: 3,
				Name:   "Item Three",
			},
			{
				ItemId: 4,
				Name:   "Item Four",
			},
		}

		info := getModelInfo(items)

		builder := newDatumBuilder(info, reflect.ValueOf(items), true)
		datums, err := builder.Keys()
		assert.NoError(t, err)
		assert.NotEmpty(t, datums)

		verify, err := builder.Verify()
		assert.NoError(t, err)
		assert.NotEmpty(t, verify)

		reader := newDatumReader(info)
		for k, v := range datums {
			value, err := reader.Read([]byte(k), v)
			assert.NoError(t, err)
			assert.NotEmpty(t, value)
		}
	})
}
