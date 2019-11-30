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

func TestQuery_InnerJoin(t *testing.T) {
	t.Skip("test")
	type ParentItem struct {
		ParentId  uint64 `m:"pk"`
		Name      string
		IsEnabled bool
	}

	type ChildItem struct {
		ChildId  uint64 `m:"pk"`
		ParentId uint64
		Parent   ParentItem `m:"fk:ParentId"`
		Name     string
	}

	parents := []ParentItem{
		{
			ParentId:  1,
			Name:      "Parent One",
			IsEnabled: true,
		},
		{
			ParentId:  2,
			Name:      "Parent Two",
			IsEnabled: false,
		},
	}

	children := []ChildItem{
		{
			ChildId:  1,
			ParentId: 1,
			Name:     "Child One",
		},
		{
			ChildId:  2,
			ParentId: 1,
			Name:     "Child Two",
		},
		{
			ChildId:  3,
			ParentId: 2,
			Name:     "Child Three",
		},
		{
			ChildId:  4,
			ParentId: 2,
			Name:     "Child Four",
		},
	}

	db, cleanup := NewTestDatabase(t)
	defer cleanup()

	txn, err := db.Begin()
	assert.NoError(t, err)

	err = txn.Insert(parents)
	assert.NoError(t, err)

	err = txn.Insert(children)
	assert.NoError(t, err)

	t.Run("simple", func(t *testing.T) {
		result := make([]ChildItem, 0)
		err = txn.
			Model(&result).
			InnerJoin(ParentItem{}).
			Where(Ex{
				"Parent.IsEnabled": true,
			}).
			Select(&result)
		assert.NoError(t, err)
		assert.Equal(t, []ChildItem{
			{
				ChildId:  1,
				ParentId: 1,
				Name:     "Child One",
			},
			{
				ChildId:  2,
				ParentId: 1,
				Name:     "Child Two",
			},
		}, result)
	})
}

func TestQuery_ScanCustomType(t *testing.T) {
	type Thing int
	type Item struct {
		ItemId uint64 `m:"pk"`
		Stuff  Thing
	}

	items := []Item{
		{
			ItemId: 1,
			Stuff:  5321,
		},
		{
			ItemId: 2,
			Stuff:  4213,
		},
		{
			ItemId: 3,
			Stuff:  6435,
		},
	}

	db, cleanup := NewTestDatabase(t)
	defer cleanup()

	txn, err := db.Begin()
	assert.NoError(t, err)

	err = txn.Insert(items)
	assert.NoError(t, err)

	read := make([]Item, 0)
	err = txn.Model(read).Select(&read)
	assert.NoError(t, err)
	assert.NotEmpty(t, read)
}
