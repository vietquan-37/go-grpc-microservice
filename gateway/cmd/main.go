package main

import (
	"common/discovery"
	"common/discovery/consul"
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/gateway/pkg/auth"
	"github.com/vietquan-37/gateway/pkg/config"
	"github.com/vietquan-37/gateway/pkg/order"
	"github.com/vietquan-37/gateway/pkg/order/pb"
	"github.com/vietquan-37/gateway/pkg/product"
	productpb "github.com/vietquan-37/gateway/pkg/product/pb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authpb "github.com/vietquan-37/gateway/pkg/auth/pb"
)

var interuptSignal = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

const (
	resolver = "consul"
)

func main() {
	c, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal().Err(err).Msg("fail to load config:")
	}
	ctx, stop := signal.NotifyContext(context.Background(), interuptSignal...)
	defer stop()
	registry, err := consul.NewRegistry(c.ConsulAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to consul")
	}
	instanceId := discovery.GenerateInstanceID(c.ServiceName)
	if err := registry.Register(instanceId, c.ServiceName, c.GatewayPort); err != nil {
		log.Fatal().Err(err).Msg("failed to register service")
	}
	err = consul.RegisterConsulResolver(registry.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register consul resolver")
	}
	go func() {
		for {
			if err := registry.HealthCheck(instanceId); err != nil {
				log.Fatal().Err(err).Msg("failed to health check service")
			}
			time.Sleep(1 * time.Second)
		}

	}()

	authClient, err := auth.InitAuthClient(c.AuthServiceName, resolver)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to init auth client")
	}
	productClient, err := product.InitProductClient(c.ProductServiceName, resolver)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to init product client")
	}
	orderClient, err := order.InitOrderClient(c.OrderServiceName, resolver)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to init order client")

	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	mux := runtime.NewServeMux(
		jsonOption,
	)

	if err = authpb.RegisterAuthServiceHandlerClient(context.Background(), mux, authClient.Client); err != nil {
		log.Fatal().Err(err).Msg("fail to register auth client: ")
	}
	if err = productpb.RegisterProductServiceHandlerClient(context.Background(), mux, productClient.Client); err != nil {
		log.Fatal().Err(err).Msg("fail to register product client:")
	}

	if err = pb.RegisterOrderServiceHandlerClient(context.Background(), mux, orderClient.Client); err != nil {
		log.Fatal().Err(err).Msg("fail to register order client:")
	}
	httpServer := &http.Server{
		Handler: mux,
		Addr:    c.GatewayPort,
	}
	waitGroup, ctx := errgroup.WithContext(ctx)
	waitGroup.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("fail to start http server")
			return err
		}
		return nil
	})
	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("shutting down http server")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("fail to shutdown http server")
			return err
		}
		log.Info().Msg("http server shutdown successfully")
		return nil
	})
	if err := waitGroup.Wait(); err != nil {
		log.Error().Err(err).Msg("fail to wait for http server")
	}
	// graceful shutdown
}
