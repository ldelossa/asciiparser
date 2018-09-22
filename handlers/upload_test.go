package handlers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	h "github.com/ldelossa/asciiparser/handlers"
	r "github.com/ldelossa/asciiparser/internal/resourcesV1"
	"github.com/stretchr/testify/assert"
)

func setupFileSizeConstrant(t *testing.T) func() {
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		err := os.Mkdir("testdata", os.ModeDir)
		if err != nil {
			t.Fatalf("failed to write testdata directory: %s", err.Error())
		}
	}

	// write 10mb file to ./testdata directory
	f1, err := os.Create("testdata/10mbfile")
	if err != nil {
		t.Fatalf("could not write test data to ./testdata directory: %s", err.Error())
	}

	f1.Truncate(1e7)

	// write 11mb file to ./testdata directory
	f2, err := os.Create("testdata/11mbfile")
	if err != nil {
		t.Fatalf("could not write test data to ./testdata directory: %s", err.Error())
	}

	f2.Truncate(1.1e7)

	return func() {
		err = os.Remove(f1.Name())
		if err != nil {
			t.Logf("failed to delete 10mb test file: %s", err.Error())
		}
		err = os.Remove(f2.Name())
		if err != nil {
			t.Logf("failed to delete 11mb test file: %s", err.Error())
		}

	}
}

var TestFileSizeTable = []struct {
	path string
	code int
}{
	{
		path: "testdata/11mbfile",
		code: http.StatusBadRequest,
	},
	{
		path: "testdata/10mbfile",
		code: http.StatusOK,
	},
}

func TestFileSizeConstrant(t *testing.T) {
	teardown := setupFileSizeConstrant(t)
	defer teardown()

	handler := h.Upload()

	for _, tt := range TestFileSizeTable {
		f, err := ioutil.ReadFile(tt.path)
		if err != nil {
			t.Fatalf("could not read file %s: %s", tt.path, err.Error())
		}

		// create request
		req := httptest.NewRequest("POST", "/api/v1/upload", bytes.NewBuffer(f))

		// create response recorder
		rr := httptest.NewRecorder()

		// call our handler directly
		handler.ServeHTTP(rr, req)

		// assert status
		assert.EqualValues(t, tt.code, rr.Code)
	}
}

var TestUploadTable = []struct {
	contents []byte
	wC       int64
	wM       map[string]int
	prefix   string
}{
	{
		contents: []byte("A sentence with no repeated words"),
		wC:       6,
		wM: map[string]int{
			"A":        1,
			"sentence": 1,
			"with":     1,
			"no":       1,
			"repeated": 1,
			"words":    1,
		},
	},
	{
		contents: []byte("A sentence with two repeated words words"),
		wC:       7,
		wM: map[string]int{
			"A":        1,
			"sentence": 1,
			"with":     1,
			"two":      1,
			"repeated": 1,
			"words":    2,
		},
	},
	{
		contents: []byte(`A sentence with three repeated words words words and a 
		new line`),
		wC: 12,
		wM: map[string]int{
			"A":        1,
			"sentence": 1,
			"with":     1,
			"three":    1,
			"repeated": 1,
			"words":    3,
			"and":      1,
			"a":        1,
			"new":      1,
			"line":     1,
		},
	},
	{
		contents: []byte("This is a very blue sentence with blueberry bluegrass and bluebell"),
		wC:       11,
		wM: map[string]int{
			"This":     1,
			"is":       1,
			"a":        1,
			"very":     1,
			"sentence": 1,
			"with":     1,
			"and":      1,
		},
	},
}

func TestUpload(t *testing.T) {
	for _, tt := range TestUploadTable {
		// create request object
		req := httptest.NewRequest("POST", "/api/v1/upload", bytes.NewBuffer(tt.contents))
		q := req.URL.Query()
		q.Add("exclude_prefix", tt.prefix)
		req.URL.RawQuery = q.Encode()

		// create response recorder
		rr := httptest.NewRecorder()
		handler := h.Upload()

		// call our handler directly
		handler.ServeHTTP(rr, req)

		// assert 200
		assert.EqualValues(t, http.StatusOK, rr.Code)

		// deserialize body
		var resp r.UploadResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		if err != nil {
			t.Fatalf("failed to deserialize response: %s", err.Error())
		}

		// assert on recevied UploadResponse
		assert.Equal(t, tt.wC, resp.WC)
		assert.Equal(t, tt.wM, resp.OC)

	}

}
