package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-zoo/bone"
	"github.com/google/uuid"
	"github.com/ldelossa/asciiparser"
	"github.com/ldelossa/asciiparser/handlers"
	r "github.com/ldelossa/asciiparser/internal/resourcesV1"
	"github.com/stretchr/testify/assert"
)

var TestGetUploadTable = []struct {
	res *r.UploadResponse
}{
	{
		res: &r.UploadResponse{
			UUID: uuid.New().String(),
			Size: 8000,
			WC:   4,
			OC: map[string]int{
				"This": 1,
				"is":   1,
				"a":    1,
				"test": 1,
			},
		},
	},
	{
		res: &r.UploadResponse{
			UUID: uuid.New().String(),
			Size: 9000,
			WC:   5,
			OC: map[string]int{
				"This":  1,
				"is":    1,
				"a":     1,
				"test":  1,
				"again": 1,
			},
		},
	},
}

func TestGetUpload(t *testing.T) {
	// create our store, handler, and a bone mux. Bone mux is necessary
	// since we use the bone.GetValue() method in our GetUpload handler.
	store := asciiparser.NewInMemStore()
	handler := handlers.GetUpload(store)
	mux := bone.New()
	mux.GetFunc("/api/v1/uploads/:id", handler)
	mux.GetFunc("/api/v1/uploads", handler)

	for _, tt := range TestGetUploadTable {
		// add resource to our store
		store.StoreV1(tt.res)

		// create our request
		path := fmt.Sprintf("/api/v1/uploads/%s", tt.res.UUID)
		req := httptest.NewRequest("GET", path, nil)

		// create response recorder
		rr := httptest.NewRecorder()

		// call our mux
		mux.ServeHTTP(rr, req)

		// assert on response
		assert.EqualValues(t, http.StatusOK, rr.Code)

		var resp r.UploadResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		if err != nil {
			t.Fatalf("failed to deserialize handler's respose: %s", err.Error())
		}

		assert.Equal(t, tt.res.UUID, resp.UUID)
		assert.Equal(t, tt.res.Size, resp.Size)
		assert.Equal(t, tt.res.OC, resp.OC)
		assert.Equal(t, tt.res.WC, resp.WC)

	}

	// test getting all uploads
	req := httptest.NewRequest("GET", "/api/v1/uploads", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusOK, rr.Code)

	// deserialize response
	var list []*r.UploadResponse
	err := json.NewDecoder(rr.Body).Decode(&list)
	if err != nil {
		t.Fatalf("failed to deserialize handler's response: %s", err.Error())
	}

	assert.Len(t, list, len(TestGetUploadTable))
}
