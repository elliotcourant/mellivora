package mellivora

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoahDbIntegration(t *testing.T) {
	db, cleanup := NewTestDatabase(t)
	defer cleanup()

	txn, err := db.Begin()
	assert.NoError(t, err)

	type DataNode struct {
		DataNodeId uint64 `m:"pk"`
		Address    string `m:"uq:uq_address_port"`
		Port       int32  `m:"uq:uq_address_port"`
		User       string
		Password   string
		Healthy    bool
	}

	type Shard struct {
		ShardId uint64 `m:"pk"`
		State   int
	}

	type Tenant struct {
		TenantId uint64 `m:"pk"`
		ShardId  uint64
	}

	type DataNodeShards struct {
		DataNodeShardId uint64 `m:"pk"`
		DataNodeId      uint64 `m:"uq:uq_data_node_id_shard_id"`
		ShardId         uint64 `m:"uq:uq_data_node_id_shard_id"`
		ReadOnly        bool
	}

	// Seed the database
	{
		dataNodes := []DataNode{
			{
				DataNodeId: 1,
				Address:    "127.0.0.1",
				Port:       5432,
				User:       "POSTGRES",
				Password:   "password",
				Healthy:    true,
			},
			{
				DataNodeId: 2,
				Address:    "127.0.0.1",
				Port:       5433,
				User:       "POSTGRES",
				Password:   "password",
				Healthy:    true,
			},
			{
				DataNodeId: 3,
				Address:    "127.0.0.1",
				Port:       5434,
				User:       "POSTGRES",
				Password:   "password",
				Healthy:    true,
			},
		}

		shards := []Shard{
			{
				ShardId: 1,
				State:   1,
			},
			{
				ShardId: 2,
				State:   1,
			},
			{
				ShardId: 3,
				State:   1,
			},
		}

		dataNodeShards := []DataNodeShards{
			{
				DataNodeShardId: 1,
				DataNodeId:      1,
				ShardId:         1,
				ReadOnly:        false,
			},
			{
				DataNodeShardId: 2,
				DataNodeId:      2,
				ShardId:         2,
				ReadOnly:        false,
			},
			{
				DataNodeShardId: 3,
				DataNodeId:      3,
				ShardId:         3,
				ReadOnly:        false,
			},

			{
				DataNodeShardId: 4,
				DataNodeId:      1,
				ShardId:         2,
				ReadOnly:        true,
			},
			{
				DataNodeShardId: 5,
				DataNodeId:      1,
				ShardId:         3,
				ReadOnly:        true,
			},

			{
				DataNodeShardId: 6,
				DataNodeId:      2,
				ShardId:         1,
				ReadOnly:        true,
			},
			{
				DataNodeShardId: 7,
				DataNodeId:      2,
				ShardId:         3,
				ReadOnly:        true,
			},

			{
				DataNodeShardId: 8,
				DataNodeId:      3,
				ShardId:         1,
				ReadOnly:        true,
			},
			{
				DataNodeShardId: 9,
				DataNodeId:      3,
				ShardId:         2,
				ReadOnly:        true,
			},
		}

		tenants := []Tenant{
			{
				TenantId: 1,
				ShardId:  1,
			},
			{
				TenantId: 2,
				ShardId:  1,
			},
			{
				TenantId: 3,
				ShardId:  2,
			},
			{
				TenantId: 4,
				ShardId:  2,
			},
			{
				TenantId: 5,
				ShardId:  3,
			},
			{
				TenantId: 6,
				ShardId:  3,
			},
		}

		err = txn.Insert(dataNodes)
		assert.NoError(t, err)

		err = txn.Insert(shards)
		assert.NoError(t, err)

		err = txn.Insert(dataNodeShards)
		assert.NoError(t, err)

		err = txn.Insert(tenants)
		assert.NoError(t, err)
	}

	t.Run("get data nodes for shard", func(t *testing.T) {
		dataNodeShards := make([]DataNodeShards, 0)
		err = txn.Model(dataNodeShards).Where(Ex{
			"ShardId":  2,
			"ReadOnly": true,
		}).Select(&dataNodeShards)
		assert.NoError(t, err)
		assert.NotEmpty(t, dataNodeShards)

		dataNodeIds := make([]uint64, 0)
		for _, dataNodeShard := range dataNodeShards {
			dataNodeIds = append(dataNodeIds, dataNodeShard.DataNodeId)
		}

		dataNodes := make([]DataNode, 0)
		err = txn.Model(dataNodes).Where(Ex{
			"DataNodeId": dataNodeIds,
		}).Select(&dataNodes)
		assert.NoError(t, err)
		assert.NotEmpty(t, dataNodes)
	})
}
