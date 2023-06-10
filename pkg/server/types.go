package server

type ParseRequestData struct {
	FileData []byte `json:"data"`
	URL      string `json:"url,omitempty"`
}
