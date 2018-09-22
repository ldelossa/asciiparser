package resourcesV1

type ErrSizeLimit struct {
	Size  int64  `json:"file_size"`
	Limit int    `json:"file_limit"`
	Diff  int64  `json:"difference"`
	Type  string `json:"type"`
}
