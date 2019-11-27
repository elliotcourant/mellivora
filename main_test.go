package mellivora

import (
	"github.com/elliotcourant/meles"
	"github.com/elliotcourant/timber"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"testing"
)

func NewTestDatabase(t *testing.T) (*Database, func()) {
	listener, err := net.Listen("tcp", ":")
	assert.NoError(t, err)
	assert.NotNil(t, listener)

	logger := timber.With(timber.Keys{
		"test": t.Name(),
	})

	tmpdir, err := ioutil.TempDir("", "mellivora")
	assert.NoError(t, err)

	store, err := meles.NewStore(listener, logger, meles.Options{
		Directory: tmpdir,
		Peers:     make([]string, 0),
	})
	assert.NoError(t, err)

	err = store.Start()
	assert.NoError(t, err)

	return &Database{
			store:  store,
			logger: logger,
		}, func() {
			store.Stop()
			listener.Close()
			os.RemoveAll(tmpdir)
		}
}

func NewCluster(t *testing.T, numberOfPeers int) ([]*Database, func()) {
	logger := timber.With(timber.Keys{
		"test": t.Name(),
	})

	listeners := make([]net.Listener, numberOfPeers)
	peers := make([]string, numberOfPeers)
	dirs := make([]string, numberOfPeers)

	for i := 0; i < numberOfPeers; i++ {
		{ // Listener
			listener, err := net.Listen("tcp", ":")
			assert.NoError(t, err)
			assert.NotNil(t, listener)
			listeners[i] = listener
			peers[i] = listener.Addr().String()
		}

		{ // Temp directory
			tmpdir, err := ioutil.TempDir("", "mellivora")
			assert.NoError(t, err)
			dirs[i] = tmpdir
		}
	}

	nodes := make([]*Database, numberOfPeers)
	for i := 0; i < numberOfPeers; i++ {
		log := logger.With(timber.Keys{
			"peer": i,
		})
		store, err := meles.NewStore(
			listeners[i],
			log,
			meles.Options{
				Directory: dirs[i],
				Peers:     peers,
			},
		)
		assert.NoError(t, err)

		nodes[i] = &Database{
			store:  store,
			logger: log,
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(numberOfPeers)
	for i := 0; i < numberOfPeers; i++ {
		go func(i int) {
			defer wg.Done()
			nodes[i].store.Start()
		}(i)
	}
	wg.Wait()

	return nodes, func() {
		for i := 0; i < numberOfPeers; i++ {
			nodes[i].store.Stop()
			listeners[i].Close()
			os.RemoveAll(dirs[i])
		}
	}
}
