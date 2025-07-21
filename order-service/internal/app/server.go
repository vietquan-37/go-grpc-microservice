package app

import (
	commonclient "common/client"
	"common/discovery"
	"common/discovery/consul"
	"common/interceptor"
	"common/kafka/producer"
	kafka_retry "common/kafka/retry"
	"common/loggers"
	"common/mtdt"
	"common/routes"
	"common/timeout"
	"common/validate"
	"context"
	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/order-service/pkg/client"
	"github.com/vietquan-37/order-service/pkg/config"
	"github.com/vietquan-37/order-service/pkg/consumer"
	"github.com/vietquan-37/order-service/pkg/db"
	"github.com/vietquan-37/order-service/pkg/handler"
	"github.com/vietquan-37/order-service/pkg/pb"
	"github.com/vietquan-37/order-service/pkg/repository"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"net"
	"sync"
	"time"
)

type ServiceDependencies struct {
	OrderHandler pb.OrderServiceServer
}

func (s *Server) setupDependencies() (*ServiceDependencies, error) {
	repo := repository.NewOrderRepo(s.db)

	productClient, paymentClient, authClient, err := s.setupExternalService()
	if err != nil {
		return nil, err
	}
	h := handler.NewOrderHandler(productClient, paymentClient, authClient, repo)
	return &ServiceDependencies{
		h,
	}, nil
}
func (s *Server) setupExternalService() (*client.ProductClient, *client.PaymentClient, *commonclient.AuthClient, error) {
	productClient, err := client.InitProductClient(s.config.ProductServiceName)
	if err != nil {
		return nil, nil, nil, err
	}
	paymentClient, err := client.InitPaymentClient(s.config.PaymentServiceName)
	if err != nil {
		return nil, nil, nil, err
	}
	authClient, err := commonclient.InitAuthClient(s.config.AuthServiceName)
	s.authClient = authClient
	if err != nil {
		return nil, nil, nil, err
	}

	return productClient, paymentClient, authClient, nil
}

type Server struct {
	config     *config.Config
	db         *gorm.DB
	registry   *consul.Registry
	authClient *commonclient.AuthClient
	instanceId string
	producer   *producer.Producer
	consumer   *kafka_retry.ConsumerWithRetry
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	grpcServer *grpc.Server
	healthSrv  *health.Server
}

func newService() *Server {
	c, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		config: c,
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}
func (s *Server) initialize() error {
	if err := s.setupDatabase(); err != nil {
		return err
	}
	if err := s.setupServiceRegistry(); err != nil {
		return err
	}
	s.setUpProducer()
	s.setupConsumer()
	return nil
}
func (s *Server) setupDatabase() error {
	s.db = db.DbConn(s.config.DbSource)
	return nil
}
func (s *Server) setupServiceRegistry() error {
	registry, err := consul.NewRegistry(s.config.ConsulAddr)
	if err != nil {
		return err
	}
	s.registry = registry
	s.instanceId = discovery.GenerateInstanceID(s.config.ServiceName)
	if err := s.registry.Register(s.instanceId, s.config.ServiceName, s.config.GrpcServerAddress); err != nil {
		return err
	}
	return consul.RegisterConsulResolver(s.registry.Client)
}
func (s *Server) setupGRPCServer() error {
	dependencies, err := s.setupDependencies()
	if err != nil {
		return err
	}
	interceptors := s.setupInterceptor()
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)
	pb.RegisterOrderServiceServer(s.grpcServer, dependencies.OrderHandler)
	reflection.Register(s.grpcServer)
	s.setupHealthCheck()
	return nil
}
func (s *Server) setupInterceptor() []grpc.UnaryServerInterceptor {
	validateInterceptor, err := validate.NewValidationInterceptor()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create validator interceptor")
	}
	roles := routes.AccessibleRoles
	authInterceptor := interceptor.NewAuthInterceptor(s.authClient, roles())
	durations := time.Second * 6
	return []grpc.UnaryServerInterceptor{
		loggers.GrpcLoggerInterceptor,
		mtdt.ForwardMetadataUnaryServerInterceptor(),
		authInterceptor.UnaryAuthInterceptor(),
		validateInterceptor.ValidateInterceptor(),
		timeout.UnaryTimeoutInterceptor(durations),
	}
}

func (s *Server) setupHealthCheck() {
	s.healthSrv = health.NewServer()
	healthpb.RegisterHealthServer(s.grpcServer, s.healthSrv)
	s.healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
}
func (s *Server) startHealthCheck() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				log.Info().Msg("Health check stopped")
				return
			case <-ticker.C:
				if err := s.registry.HealthCheck(s.instanceId); err != nil {
					log.Error().Err(err).Msg("Health check failed")
				}
			}
		}
	}()
}
func (s *Server) gracefulShutdown() {
	log.Info().Msg("Shutting down server...")

	if s.healthSrv != nil {
		s.healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	}

	if s.grpcServer != nil {
		log.Info().Msg("Stopping gRPC server...")
		s.grpcServer.GracefulStop()
	}

	if s.cancel != nil {
		s.cancel()
	}
	log.Info().Msg("Waiting for all goroutines to finish...")
	s.wg.Wait()
	log.Info().Msg("Server stopped gracefully")
}

func (s *Server) cleanup() {
	if s.registry != nil && s.instanceId != "" {
		s.registry.Deregister(s.instanceId, s.config.ServiceName)
	}
}
func (s *Server) start() error {
	lis, err := net.Listen("tcp", s.config.GrpcServerAddress)
	if err != nil {
		return err
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Info().Msg("Starting Kafka consumer...")
		if err := s.consumer.Start(s.ctx); err != nil {
			log.Error().Err(err).Msg("Kafka consumer stopped with error")
		}
		log.Info().Msg("Kafka consumer stopped")
	}()

	go func() {
		log.Info().Msgf("Starting gRPC server at %s", lis.Addr().String())
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	log.Info().Msg("gRPC server started")
	return nil
}
func (s *Server) setUpProducer() {
	s.producer = producer.NewProducer(s.config.BrokerAddr...)
}

func (s *Server) setupConsumer() {
	repo := repository.NewOrderRepo(s.db)
	productConsumer := consumer.NewOrderConsumer(repo, s.config, s.producer)
	s.consumer = kafka_retry.NewConsumerWithRetry(
		s.config.BrokerAddr,
		s.config.OrderTopic,
		s.config.GroupId,
		productConsumer,
		s.config.MaxRetries,
		func() backoff.BackOff {
			return backoff.NewExponentialBackOff(
				backoff.WithMaxElapsedTime(20 * time.Second),
			)
		},
		s.config.WorkerCount,
	)
}
