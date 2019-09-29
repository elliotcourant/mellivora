package mellivora

import (
	"github.com/elliotcourant/meles"
	"github.com/elliotcourant/timber"
	"net"
)

type Store struct {
	listener net.Listener
	db       *meles.Store
}

type Options struct {
	Directory     string
	ListenAddress string
	Peers         []string
	Logger        timber.Logger
}

var DefaultOptions = Options{
	Directory:     "data",
	ListenAddress: ":",
	Peers:         make([]string, 0),
	Logger:        timber.New(),
}

// Open opens or creates a new data store.
func Open(options Options) (*Store, error) {
	listener, err := net.Listen("tcp", options.ListenAddress)
	if err != nil {
		return nil, err
	}

	store := &Store{
		listener: listener,
		db:       nil,
	}

	db, err := meles.NewStore(listener, options.Logger, meles.Options{
		Directory: options.Directory,
		Peers:     options.Peers,
	})
	if err != nil {
		return nil, err
	}
	store.db = db
	if err := store.db.Start(); err != nil {
		return nil, err
	}
	return store, nil
}

// Meles returns the underlying Meles store.
func (s *Store) Meles() *meles.Store {
	return s.db
}

// Close closes and stops the data store.
func (s *Store) Close() error {
	return s.db.Stop()
}
