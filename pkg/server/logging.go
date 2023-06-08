package server

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type LoggingService struct {
	logger zerolog.Logger
	next   Service
}

func NewLoggingService(next Service) Service {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	return &LoggingService{
		next:   next,
		logger: logger,
	}
}

func (s *LoggingService) PostStateFile(stateFile []byte, ctx context.Context) (parsedFile []byte, err error) {
	// defer is called when this function returns
	// this defer enables the named values to be used
	defer func() {
		if err != nil {
			s.logger.Error().Err(err).Msg("")
		} else {
			s.logger.Info().Msg("parsed state file")
		}
	}()
	return s.next.PostStateFile(stateFile, ctx)
}
