package mellivora

import (
	"github.com/elliotcourant/meles"
	"github.com/elliotcourant/timber"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"os"
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

	return &Database{store: store}, func() {
		store.Stop()
		listener.Close()
		os.RemoveAll(tmpdir)
	}
}
