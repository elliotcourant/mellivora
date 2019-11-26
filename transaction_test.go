package mellivora

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction_Insert(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		db, cleanup := NewTestDatabase(t)
		defer cleanup()

		txn, err := db.Begin()
		assert.NoError(t, err)

		type DataNode struct {
			DataNodeId uint64 `m:"pk,serial"`
			Address    string `m:"unique:uq_address_port"`
			Port       int32  `m:"unique:uq_address_port"`
			User       string
			Password   string
			Healthy    bool
		}

		dataNodes := []*DataNode{
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

		err = txn.Insert(dataNodes)
		assert.NoError(t, err)

		err = txn.Commit()
		assert.NoError(t, err)
	})
}