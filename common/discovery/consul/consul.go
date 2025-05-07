package consul

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
	"net"
	"strconv"
	"strings"
)

func resolveHostPort(hostPort string, mode string) (string, int, error) {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 || parts[1] == "" {
		return "", 0, fmt.Errorf("invalid host port format: %s", hostPort)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %v", err)
	}

	switch mode {
	case "local":
		return "127.0.0.1", port, nil

	case "production":
		var ips []string
		addresses, err := net.InterfaceAddrs()
		if err != nil {
			return "", 0, fmt.Errorf("error getting interface addresses: %v", err)
		}

		for _, addr := range addresses {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {

					ips = append(ips, ipnet.IP.String())
				}
			}
		}

		if len(ips) == 0 {
			return "", 0, fmt.Errorf("could not find a suitable non-loopback IP address")
		}
		log.Info().Msgf("ips : %v", ips)

		return ips[0], port, nil

	default:
		return "", 0, fmt.Errorf("unsupported mode: %s", mode)
	}
}

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
func (r *Registry) Register(instanceID, serviceName, hostPort, mode string) error {
	host, port, err := resolveHostPort(hostPort, mode)
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
