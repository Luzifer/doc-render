// Package mem implements a storage backend to hold the templates
// inside memory for testing purposes
package mem

import (
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/Luzifer/doc-render/pkg/persist"
)

type (
	// Backend implements the persist.Backend interface for Memory storage
	Backend struct {
		store map[string][]byte
		lock  sync.RWMutex
	}
)

var _ persist.Backend = (*Backend)(nil)

// New creates a new redis persistence backend
func New() *Backend {
	return &Backend{
		store: map[string][]byte{},
	}
}

// Get retrieves the JSON encoded template by its content-hash
func (b *Backend) Get(uid string) (templateJSON []byte, err error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.store[uid], nil
}

// Store takes the JSON encoded template, stores it and returns the
// content-hash as uid and optionally an error
func (b *Backend) Store(templateJSON []byte) (uid string, err error) {
	uid = b.contentHash(templateJSON)

	b.lock.Lock()
	defer b.lock.Unlock()

	b.store[uid] = templateJSON

	return uid, nil
}

func (*Backend) contentHash(templateJSON []byte) string {
	return fmt.Sprintf("sha256:%x", sha256.Sum256(templateJSON))
}
