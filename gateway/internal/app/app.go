package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var interruptSignal = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

type App struct {
	server *Server
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignal...)
	defer stop()
	a.server = newServer()

	if err := a.server.initialize(); err != nil {
		return err
	}
	defer a.server.cleanup()

	if err := a.server.setupHTTPServer(); err != nil {
		return err
	}

	go a.server.startHealthCheck()

	return a.server.start(ctx)
}
