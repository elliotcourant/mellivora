package mellivora

import (
	"fmt"
	"github.com/elliotcourant/meles"
	"reflect"
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
	info := getModelInfo(model)
	builder := newDatumBuilder(info, reflect.ValueOf(model))

	datums, err := builder.Keys()
	if err != nil {
		return err
	}

	verify, err := builder.Verify()
	if err != nil {
		return err
	}

	for verifyKey, canExist := range verify {
		_, ok, err := txn.tx.MustGet([]byte(verifyKey))
		if err != nil {
			return err
		}
		if ok && !canExist {
			return fmt.Errorf("an item with key [%s] already exists, cannot overwrite", verifyKey)
		}
	}

	for key, value := range datums {
		if err := txn.tx.Set([]byte(key), value); err != nil {
			return err
		}
	}

	return nil
}
