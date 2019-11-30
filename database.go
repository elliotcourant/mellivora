package mellivora

import (
	"github.com/elliotcourant/meles"
	"github.com/elliotcourant/timber"
)

type Database struct {
	store  *meles.Store
	logger timber.Logger
}

func NewDatabase(store *meles.Store, logger timber.Logger) *Database {
	return &Database{
		store:  store,
		logger: logger,
	}
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
