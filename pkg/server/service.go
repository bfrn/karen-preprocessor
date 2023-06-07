package server

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/bfrn/karen-preprocessor/pkg/preprocessor"
)

type Service interface {
	PostStateFile([]byte, context.Context) ([]byte, error)
}

type PreprocessorService struct {
}

func NewPreprocessorService() Service {
	return &PreprocessorService{}
}

func (p *PreprocessorService) PostStateFile(stateFile []byte, ctx context.Context) ([]byte, error) {
	if len(stateFile) == 0 {
		return nil, errors.New("empty state file")
	}
	parsedModel, err := preprocessor.ParseStateFile(stateFile)
	if err != nil {
		return nil, err
	}

	return json.Marshal(parsedModel)
}
