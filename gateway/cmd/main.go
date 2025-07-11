package main

import (
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/gateway/internal/app"
)

func main() {
	application := app.NewApp()
	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run application")
	}
}
