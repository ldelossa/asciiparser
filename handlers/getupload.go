package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-zoo/bone"
	ap "github.com/ldelossa/asciiparser"
	"github.com/ldelossa/asciiparser/internal/resourcesV1"
)

func GetUpload(store ap.Storer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// look for UUID of upload, we use the bone method here
		// since bone muxer is responsible for capturing url parameters
		id := bone.GetValue(r, "id")

		// block to handle a request for a single ID
		if id != "" {

			res, ok, err := store.GetV1(id)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// return upload resource
			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			return
		}

		// no ID was specified, return all info
		var list []*resourcesV1.UploadResponse
		list = store.GetAllV1()

		err := json.NewEncoder(w).Encode(list)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}
