package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"

	"github.com/bfrn/karen-preprocessor/pkg/preprocessor"
)

type Service interface {
	PostParseFile([]byte, context.Context) ([]byte, error)
}

type PreprocessorService struct {
}

func NewPreprocessorService() Service {
	return &PreprocessorService{}
}

func (p *PreprocessorService) PostParseFile(data []byte, ctx context.Context) ([]byte, error) {
	var parseRequestData *ParseRequestData
	var err error
	var parsedModel map[string]preprocessor.Node

	err = json.Unmarshal(data, &parseRequestData)
	if err != nil {
		return nil, err
	}
	if len(parseRequestData.FileData) == 0 {
		return nil, errors.New("empty file")
	}

	switch parseRequestData.FileType {
	case State:
		parsedModel, err = preprocessor.ParseStateFile([]byte(parseRequestData.FileData))
		if err != nil {
			return nil, err
		}
	case Plan:
		var u *url.URL
		u, err = url.Parse(parseRequestData.URL)
		if err != nil {
			return nil, err
		}
		parsedModel, err = preprocessor.ParsePlanFile([]byte(parseRequestData.FileData), u.Host, u.Path)
		if err != nil {
			return nil, err
		}
	case Karen:
		return data, nil
	default:
		return nil, errors.New("unknown file format")
	}

	return json.Marshal(parsedModel)
}
