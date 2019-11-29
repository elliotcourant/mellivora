package mellivora

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery_Where(t *testing.T) {
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
	}

	db, cleanup := NewTestDatabase(t)
	defer cleanup()

	txn, err := db.Begin()
	assert.NoError(t, err)

	err = txn.Insert(items)
	assert.NoError(t, err)

	t.Run("filter by name", func(t *testing.T) {
		result := make([]Item, 0)
		err = txn.Model(result).Where(Ex{
			"Name": "Item Two",
		}).Select(&result)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, []Item{
			{
				ItemId: 2,
				Name:   "Item Two",
			},
		}, result)
	})

	t.Run("filter by id", func(t *testing.T) {
		result := make([]Item, 0)
		err = txn.Model(result).Where(Ex{
			"ItemId": 3,
		}).Select(&result)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, []Item{
			{
				ItemId: 3,
				Name:   "Item Three",
			},
		}, result)
	})

	t.Run("filter by id into struct", func(t *testing.T) {
		result := Item{}
		err = txn.Model(result).Where(Ex{
			"ItemId": 3,
		}).Select(&result)
		assert.NoError(t, err)
		assert.Equal(t, Item{
			ItemId: 3,
			Name:   "Item Three",
		}, result)
	})
}
