package main

import (
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/email-service/internal/app"
)

func main() {
	application := app.NewApp()

	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run email service")
	}
}
