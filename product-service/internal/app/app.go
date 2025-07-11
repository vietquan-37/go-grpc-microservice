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

	a.server = newService()
	if err := a.server.initialize(); err != nil {
		return err
	}
	defer a.server.cleanup()
	if err := a.server.setupGrpcServer(); err != nil {
		return err
	}
	a.server.startHealthCheck()
	if err := a.server.Start(); err != nil {
		return err
	}
	<-ctx.Done()
	a.server.gracefulShutdown()
	return nil

}
