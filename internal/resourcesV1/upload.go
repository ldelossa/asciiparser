package resourcesV1

type UploadResponse struct {
	Size int64          `json:"size"`
	UUID string         `json:"id"`
	WC   int64          `json:"word_count"`
	OC   map[string]int `json:"occurences"`
}
