package asciiparser_test

import (
	"testing"

	"github.com/google/uuid"
	ap "github.com/ldelossa/asciiparser"
	"github.com/ldelossa/asciiparser/internal/resourcesV1"
	"github.com/stretchr/testify/assert"
)

var InMemStoreTestTable = []struct {
	res *resourcesV1.UploadResponse
}{
	{
		res: &resourcesV1.UploadResponse{
			UUID: uuid.New().String(),
			WC:   5,
			OC: map[string]int{
				"One":   1,
				"Two":   1,
				"Three": 1,
				"Four":  1,
				"Five":  1,
			},
		},
	},
	{
		res: &resourcesV1.UploadResponse{
			UUID: uuid.New().String(),
			WC:   6,
			OC: map[string]int{
				"One":   1,
				"Two":   1,
				"Three": 1,
				"Four":  1,
				"Five":  1,
				"Six":   1,
			},
		},
	},
}

func TestInMemStore(t *testing.T) {
	store := ap.NewInMemStore()

	for _, tt := range InMemStoreTestTable {
		// store object
		err := store.StoreV1(tt.res)
		if err != nil {
			t.Fatalf("failed to store resource: %s", err.Error())
		}

		// get from store
		res, ok, err := store.GetV1(tt.res.UUID)
		if err != nil {
			t.Fatalf("failed to get resource UUID: %s", err.Error())
		}
		if !ok {
			t.Fatalf("uuid not found")
		}

		// assert on recevied record
		assert.Equal(t, tt.res.UUID, res.UUID)
		assert.Equal(t, tt.res.WC, res.WC)
		assert.Equal(t, tt.res.OC, res.OC)
	}

	// attempt to get all objects
	list := store.GetAllV1()

	// assert length
	assert.Len(t, list, len(InMemStoreTestTable))

}
