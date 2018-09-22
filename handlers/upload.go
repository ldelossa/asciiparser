package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	res "github.com/ldelossa/asciiparser/internal/resourcesV1"
)

const (
	sizeLimitBytes = 10000000
	excludePrefix  = "blue"
)

// Upload is our handler for uploading an ASCII file to. We use a closure syntax incase future dependencies
// must be plumbed into handler.
func Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// this is a POST hold handler
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// confirm message does not exceed 10mb. ContentLength is size of body in bytes
		if r.ContentLength > sizeLimitBytes {
			log.Printf("received content exceeding 10mb")
			w.WriteHeader(http.StatusBadRequest)

			// write response body
			err := json.NewEncoder(w).Encode(res.ErrSizeLimit{
				Limit: sizeLimitBytes,
				Size:  r.ContentLength,
				Type:  "SizeLimitExceeded",
				Diff:  r.ContentLength - sizeLimitBytes,
			})
			if err != nil {
				log.Printf("failed to serialize ErrSizeLimit struct")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		// create a map we will use to hold word counts
		// create word count variable
		var wm map[string]int = map[string]int{}
		var wc int64
		prefix := []byte(excludePrefix)

		// create bufio Scanner from request body
		s := bufio.NewScanner(r.Body)
		// set the scanner's split function, we utilize the default implementation's ScanWords method
		s.Split(bufio.ScanWords)

		// start parsing our words
		for s.Scan() {
			// increate our word count
			wc++

			// first check that prefix is not not empty (empty string byte array is len 0)
			// next use byte equality method to compare current word's byte array prefix length
			// to our prefix byte array. If they match do not add current word to word map
			if len(prefix) > 0 && bytes.Equal(s.Bytes()[:len(prefix)], prefix) {
				continue
			}

			// add current word to map
			wm[s.Text()] = wm[s.Text()] + 1
		}

		// return word count and word map
		err := json.NewEncoder(w).Encode(res.UploadResponse{
			Size: r.ContentLength,
			UUID: uuid.New().String(),
			WC:   wc,
			OC:   wm,
		})
		if err != nil {
			log.Printf("failed to serialize UploadResponse")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// implicity 200 if no other header is written
		return

	}
}
