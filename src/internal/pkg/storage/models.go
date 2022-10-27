package storage

import "time"

type UploadFileResponse struct {
	Hash      string    `json:"hash,omitempty"`
	Size      int64     `json:"size,omitempty"`
	TimeStamp time.Time `json:"timestamp,omitempty"`
}

type AllFilesResponse struct {
	Results []string `json:"results,omitempty"`
	Count   int      `json:"count,omitempty"`
}
