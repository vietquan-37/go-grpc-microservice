package consul

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
	"net"

	"strconv"
	"strings"
)

type Registry struct {
	Client *consul.Client
}

func resolveHostPort(hostPort string, resolve bool) (string, int, error) {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid host port %s", hostPort)
	}
	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, err
	}
	if !resolve || net.ParseIP(host) != nil || host == "localhost" {
		return host, port, nil
	}

	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		return "", 0, fmt.Errorf("failed to resolve host %s: %v", host, err)
	}
	realHost := ips[0].String()
	return realHost, port, nil
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
func (r *Registry) Register(instanceID, serviceName, hostPort string, resolve bool) error {
	host, port, err := resolveHostPort(hostPort, resolve)
	if err != nil {
		return err
	}
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
