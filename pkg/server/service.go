package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"

	"github.com/bfrn/karen-preprocessor/pkg/preprocessor"
)

type Service interface {
	PostParseFile(*ParseRequestData, context.Context) ([]byte, error)
}

type PreprocessorService struct {
}

func NewPreprocessorService() Service {
	return &PreprocessorService{}
}

func (p *PreprocessorService) PostParseFile(data *ParseRequestData, ctx context.Context) ([]byte, error) {
	if len(data.FileData) == 0 {
		return nil, errors.New("empty state file")
	}
	parsedModel, err := preprocessor.ParseStateFile(data.FileData)
	if err == nil {
		return json.Marshal(parsedModel)
	}
	u, err := url.Parse(data.URL)
	if err != nil {
		return nil, errors.New("plan file require an valid URL to the entry terraform file, e.g. 'https://github.com/example/main.tf'")
	}
	parsedModel, err = preprocessor.ParsePlanFile(data.FileData, u.Host, u.Path)
	if err == nil {
		return json.Marshal(parsedModel)
	}

	return nil, errors.New("invalid file content: provide a valid state or plan file")
}
