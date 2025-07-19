package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	if err := a.server.setupGRPCServer(); err != nil {
		return err
	}

	a.server.startHealthCheck()

	if err := a.server.start(); err != nil {
		return err
	}

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	a.server.gracefulShutdown(shutdownCtx)
	return nil
}
