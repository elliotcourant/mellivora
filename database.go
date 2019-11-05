package mellivora

import (
	"github.com/elliotcourant/meles"
)

type Database struct {
	store *meles.Store
}

func Open() *Database {
	return nil
}

func (db *Database) Begin() (*Transaction, error) {
	storeTxn, err := db.store.Begin()
	if err != nil {
		return nil, err
	}

	return &Transaction{
		db: db,
		tx: storeTxn,
	}, err
}
