package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

type Resolver struct {
	serviceName string
	client      *api.Client
	cc          resolver.ClientConn
	mu          sync.Mutex
	stopCh      chan struct{}
}

func NewConsulResolver(serviceName string, cc resolver.ClientConn, client *api.Client) *Resolver {
	r := &Resolver{
		serviceName: serviceName,
		client:      client,
		cc:          cc,
		stopCh:      make(chan struct{}),
	}
	go r.watchConsul()
	return r
}

func (r *Resolver) resolve() {
	r.mu.Lock()
	defer r.mu.Unlock()

	services, _, err := r.client.Health().Service(r.serviceName, "", true, nil)

	if err != nil {
		log.Error().Err(err).Msg("Consul resolver get services failed")
		return
	}

	var addresses []resolver.Address
	for _, service := range services {

		addr := fmt.Sprintf("%s:%d", service.Service.Port, service.Service.Port)
		addresses = append(addresses, resolver.Address{Addr: addr})
	}

	if err := r.cc.UpdateState(resolver.State{Addresses: addresses}); err != nil {
		//log.Error().Err(err).Msg("Consul resolver update state failed")
		return
	}
}

func (r *Resolver) watchConsul() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.stopCh:
			return
		case <-ticker.C:
			r.resolve()
		}
	}
}

// loadbalancing change and resolve  Immediate
func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {

	go r.resolve()
}

func (r *Resolver) Close() {
	close(r.stopCh)
}
