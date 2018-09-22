package asciiparser

import (
	"fmt"
	"sync"

	"github.com/ldelossa/asciiparser/internal/resourcesV1"
)

// Storer provides an interface for storing and retrieving
// UploadResponse resources
type Storer interface {
	StoreV1(*resourcesV1.UploadResponse) error
	GetV1(uuid string) (*resourcesV1.UploadResponse, bool, error)
	GetAllV1() []*resourcesV1.UploadResponse
}

// InMemStore creates a temporary in memory store. This can be used
// for development servers or testing in mocks
type InMemStore struct {
	// use a syncmap here - normally I would write the locking mechanisms but will keep this definition simple
	m sync.Map
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		m: sync.Map{},
	}
}

func (i *InMemStore) StoreV1(r *resourcesV1.UploadResponse) error {
	i.m.Store(r.UUID, r)
	return nil
}

func (i *InMemStore) GetV1(UUID string) (*resourcesV1.UploadResponse, bool, error) {
	r, ok := i.m.Load(UUID)
	if !ok {
		return nil, ok, fmt.Errorf("failed to find resource with UUID %s", UUID)
	}

	return r.(*resourcesV1.UploadResponse), true, nil
}

func (i *InMemStore) GetAllV1() []*resourcesV1.UploadResponse {
	var list []*resourcesV1.UploadResponse

	// func has access to list via closure
	i.m.Range(func(key interface{}, value interface{}) bool {
		list = append(list, value.(*resourcesV1.UploadResponse))
		return true
	})

	return list
}
