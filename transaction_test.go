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

		// Insert the records initially, this insert should
		// succeed as there are no records in the DB yet.
		{
			err = txn.Insert(dataNodes)
			assert.NoError(t, err)

			err = txn.Commit()
			assert.NoError(t, err)
		}

		// Insert the records again, this insert should
		// fail due to the unique constraint so make sure
		// an error is returned from the insert call.
		{
			txn, err = db.Begin()
			assert.NoError(t, err)

			err = txn.Insert(dataNodes)
			assert.Error(t, err)
		}
	})

	t.Run("cluster", func(t *testing.T) {
		cluster, cleanup := NewCluster(t, 3)
		defer cleanup()

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

		txn1, err := cluster[0].Begin()
		assert.NoError(t, err)

		err = txn1.Insert(dataNodes)
		assert.NoError(t, err)

		err = txn1.Commit()
		assert.NoError(t, err)
	})

	t.Run("conflict", func(t *testing.T) {
		cluster, cleanup := NewCluster(t, 3)
		defer cleanup()

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

		txn1, err := cluster[0].Begin()
		assert.NoError(t, err)

		txn2, err := cluster[1].Begin()
		assert.NoError(t, err)

		// Insert the records initially, this insert should
		// succeed as there are no records in the DB yet.
		{
			err = txn1.Insert(dataNodes)
			assert.NoError(t, err)
		}

		// Insert the records again, this insert should
		// fail due to the unique constraint so make sure
		// an error is returned from the insert call.
		{
			err = txn2.Insert(dataNodes)
			assert.NoError(t, err)
		}

		{
			err = txn1.Commit()
			assert.NoError(t, err)

			err = txn2.Commit()
			assert.Error(t, err)
		}
	})
}
