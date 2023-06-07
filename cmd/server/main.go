package main

import (
	"github.com/bfrn/karen-preprocessor/pkg/server"
	"github.com/rs/zerolog/log"
)

func main() {

	svc := server.NewPreprocessorService()
	svc = server.NewLoggingService(svc)

	apiServer := server.NewApiServer(svc)
	log.Fatal().Err(apiServer.Start(":3000"))
}
