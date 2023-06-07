package util

import (
	"os"

	"github.com/rs/zerolog/log"
)

func CheckErr(err error) {
	if err == nil {
		return
	}
	log.Fatal().Err(err).Msg("")
	os.Exit(1)
}
