package server

type FileType string

const (
	Karen     FileType = "karen"
	Plan      FileType = "plan"
	State     FileType = "state"
	Undefined FileType = "undefined"
)

type ParseRequestData struct {
	FileData string   `json:"data"`
	FileType FileType `json:"type"`
	URL      string   `json:"url,omitempty"`
}
