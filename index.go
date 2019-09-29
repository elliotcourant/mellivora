package mellivora

import (
	"github.com/elliotcourant/meles"
)

// IndexFunction is the actual function that us used to generate
// the bytes stored for a value.
type IndexFunction func(name string, value interface{}) ([]byte, error)

// Index is wrapper that returns the indexable encoded bytes of the provided value.
type Index struct {
	IndexFunction IndexFunction
	Unique        bool
}

func indexAdd(storer Instructions, txn *meles.Transaction, key []byte, data interface{}) error {
	panic("not implemented")
	// indexes := storer.Indexes()
	// for name, index := range indexes {
	//
	// }
}
