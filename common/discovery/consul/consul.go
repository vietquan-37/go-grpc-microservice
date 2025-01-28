package consul

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"

	"strconv"
	"strings"
)

type Registry struct {
	Client *consul.Client
}

func NewRegistry(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client}, nil

}
func (r *Registry) Register(instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid host port %s", hostPort)
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	host := parts[0]
	return r.Client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Address: host,
		Check: &consul.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}
func (r *Registry) Deregister(instanceID, serviceName string) {
	log.Info().Msgf("Deregistering %s %s", instanceID, serviceName)
	if err := r.Client.Agent().ServiceDeregister(instanceID); err != nil {
		log.Fatal().Err(err).Msg("fail deregister service: ")
	}
}

func (r *Registry) HealthCheck(instanceID string) error {
	return r.Client.Agent().UpdateTTL(
		instanceID,
		"online",
		consul.HealthPassing,
	)
}
