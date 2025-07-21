package app

import (
	"common/kafka/consumer"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/email-service/pkg/config"
	"github.com/vietquan-37/email-service/pkg/email"
	"github.com/vietquan-37/email-service/pkg/handler"

	"sync"
)

type Server struct {
	config         *config.Config
	emailService   email.IEmailService
	messageHandler *handler.MessageHandler
	consumer       *consumer.Consumer
	ctx            context.Context
	cancel         context.CancelFunc
	wg             *sync.WaitGroup
}

func newServer() *Server {
	cfg, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}

func (s *Server) initialize() error {
	if err := s.setupDependencies(); err != nil {
		return err
	}

	if err := s.setupConsumer(); err != nil {
		return err
	}

	return nil
}

func (s *Server) setupDependencies() (err error) {

	s.emailService, err = email.NewEmailService(s.config)
	if err != nil {
		return err
	}

	s.messageHandler = handler.NewMessageHandler(s.emailService)

	return nil
}

func (s *Server) setupConsumer() error {

	s.consumer = consumer.NewConsumer(
		s.config.BrokerAddr,
		s.config.Topic,
		s.config.GroupID,
		s.messageHandler.ProcessMessage,
	)

	return nil
}

func (s *Server) start() error {
	log.Info().
		Str("topic", s.config.Topic).
		Str("group_id", s.config.GroupID).
		Msg("Starting email service consumer...")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.consumer.Start(s.ctx); err != nil {
			log.Error().Err(err).Msg("Consumer stopped with error")
		}
	}()

	log.Info().Msg("Email service consumer started")
	return nil
}

func (s *Server) gracefulShutdown() {
	log.Info().Msg("Shutting down email service...")

	// Cancel context to stop consumer
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
	log.Info().Msg("Waiting for goroutines to finish")

	log.Info().Msg("Email service stopped gracefully")
}
