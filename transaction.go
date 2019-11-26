package mellivora

import (
	"github.com/elliotcourant/meles"
)

type Transaction struct {
	db *Database
	tx *meles.Transaction
}

func (txn *Transaction) Model(model interface{}) *Query {
	return &Query{
		// model: model,
		txn: txn,
	}
}

func (txn *Transaction) Commit() error {
	return txn.tx.Commit()
}

func (txn *Transaction) Rollback() error {
	return txn.tx.Rollback()
}

func (txn *Transaction) Insert(model interface{}) error {
	// info := getModelInfo(model)
	return nil
}
