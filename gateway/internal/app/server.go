package app

import (
	"common/discovery"
	"common/discovery/consul"
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/gateway/pkg/auth"
	authpb "github.com/vietquan-37/gateway/pkg/auth/pb"
	"github.com/vietquan-37/gateway/pkg/config"
	"github.com/vietquan-37/gateway/pkg/middleware"
	"github.com/vietquan-37/gateway/pkg/order"
	"github.com/vietquan-37/gateway/pkg/order/pb"
	"github.com/vietquan-37/gateway/pkg/product"
	productpb "github.com/vietquan-37/gateway/pkg/product/pb"
	"github.com/vietquan-37/gateway/pkg/ratelimiter"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"time"
)

const (
	resolver = "consul"
)

type Server struct {
	config     *config.Config
	registry   *consul.Registry
	instanceId string
	ctx        context.Context
	cancel     context.CancelFunc
	httpServer *http.Server
	clients    *ServiceClients
}

type ServiceClients struct {
	AuthClient    *auth.Client
	ProductClient *product.Client
	OrderClient   *order.Client
}

func newServer() *Server {
	cfg, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal().Err(err).Msg("fail to load config")
	}
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Server) initialize() error {
	if err := s.setupServiceRegistry(); err != nil {
		return err
	}

	if err := s.setupServiceClients(); err != nil {
		return err
	}

	return nil
}

func (s *Server) setupServiceRegistry() error {
	registry, err := consul.NewRegistry(s.config.ConsulAddr)
	if err != nil {
		return err
	}
	s.registry = registry
	s.instanceId = discovery.GenerateInstanceID(s.config.ServiceName)
	if err := s.registry.Register(s.instanceId, s.config.ServiceName, s.config.GatewayPort); err != nil {
		return err
	}
	return consul.RegisterConsulResolver(s.registry.Client)
}

func (s *Server) setupServiceClients() error {
	authClient, err := auth.InitAuthClient(s.config.AuthServiceName, resolver)
	if err != nil {
		return err
	}

	productClient, err := product.InitProductClient(s.config.ProductServiceName, resolver)
	if err != nil {
		return err
	}

	orderClient, err := order.InitOrderClient(s.config.OrderServiceName, resolver)
	if err != nil {
		return err
	}

	s.clients = &ServiceClients{
		AuthClient:    authClient,
		ProductClient: productClient,
		OrderClient:   orderClient,
	}

	return nil
}

func (s *Server) setupHTTPServer() error {
	mux, err := s.setupServeMux()
	if err != nil {
		return err
	}

	if err := s.registerServiceHandlers(mux); err != nil {
		return err
	}

	handler := s.setupMiddleware(mux)

	s.httpServer = &http.Server{
		Handler: handler,
		Addr:    s.config.GatewayPort,
	}

	return nil
}

func (s *Server) setupServeMux() (*runtime.ServeMux, error) {
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	return runtime.NewServeMux(jsonOption), nil
}

func (s *Server) registerServiceHandlers(mux *runtime.ServeMux) error {
	if err := authpb.RegisterAuthServiceHandlerClient(context.Background(), mux, s.clients.AuthClient.Client); err != nil {
		return err
	}

	if err := productpb.RegisterProductServiceHandlerClient(context.Background(), mux, s.clients.ProductClient.Client); err != nil {
		return err
	}

	if err := pb.RegisterOrderServiceHandlerClient(context.Background(), mux, s.clients.OrderClient.Client); err != nil {
		return err
	}

	return nil
}

func (s *Server) setupMiddleware(mux *runtime.ServeMux) http.Handler {
	rateLimiter := ratelimiter.NewFixedWindowRateLimiter(
		s.config.RequestPerTimeFrame,
		time.Second*5)

	rlMiddleware := middleware.NewRateLimiterMiddleware(rateLimiter)
	corsMiddleware := middleware.CorsMiddleware(rlMiddleware.RateLimitMiddleware(mux))
	return middleware.HealthCheckHandler(corsMiddleware)
}

func (s *Server) start(ctx context.Context) error {
	waitGroup, ctx := errgroup.WithContext(ctx)

	waitGroup.Go(func() error {
		log.Info().Msgf("Starting HTTP server at %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("HTTP server failed")
			return err
		}
		log.Info().Msg("http server shutdown successfully")
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		return s.gracefulShutdown()
	})

	return waitGroup.Wait()
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

func (s *Server) gracefulShutdown() error {
	log.Info().Msg("Shutting down server...")
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("fail to shutdown http server")
			return err
		}
	}

	log.Info().Msg("HTTP server shutdown successfully")
	return nil
}

func (s *Server) cleanup() {
	if s.registry != nil && s.instanceId != "" {
		s.registry.Deregister(s.instanceId, s.config.ServiceName)
	}
}
