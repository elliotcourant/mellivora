package mellivora

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery_Where(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		t.Skip("in progress")
		db, cleanup := NewTestDatabase(t)
		defer cleanup()

		txn, err := db.Begin()
		assert.NoError(t, err)

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
		}

		err = txn.Insert(items)
		assert.NoError(t, err)

		result := make([]*Item, 0)
		err = txn.Model(result).Where(Ex{
			"Name": "Item Two",
		}).Select(&result)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})
}
