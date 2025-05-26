package consul

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
	"net"
	"strconv"
)

func isPrivateIP(ipAddr string) bool {
	var privateBlocks []*net.IPNet
	for _, b := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10"} {
		if _, block, err := net.ParseCIDR(b); err == nil {
			log.Info().Msgf("block : %s", block.String())
			privateBlocks = append(privateBlocks, block)
		}
	}

	ip := net.ParseIP(ipAddr)
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}

	return false
}

func ExtractRealIP(ip string) (string, error) {
	if len(ip) > 0 && (ip != "0.0.0.0" && ip != "[::]" && ip != "::") {
		return ip, nil
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("Failed to get interface addresses, error  %s", err)
	}

	var ipAddr []byte
	for _, rawAddr := range addrs {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}

		log.Info().Msgf("Found IP: %s", ip.String())

		if ip.To4() == nil {
			continue
		}

		log.Info().Msgf("IPv4 Address: %s", ip.String())

		if !isPrivateIP(ip.String()) {
			log.Info().Msgf("Ignoring public IPv4: %s", ip.String())
			continue
		}

		log.Info().Msgf("Using private IP: %s", ip.String())
		ipAddr = ip
		break
	}

	if ipAddr == nil {
		return "", fmt.Errorf("no private IP address found, and explicit IP not provided")
	}

	return net.IP(ipAddr).String(), nil
}

func ExtractHostPortFromIP(listen string) (string, int, error) {
	host, port, err := net.SplitHostPort(listen)
	if err != nil {
		return "", 0, err
	}

	host, err = ExtractRealIP(host)
	if err != nil {
		return "", 0, err
	}

	intPort, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return "", 0, err
	}

	return host, int(intPort), nil
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
	host, port, err := ExtractHostPortFromIP(hostPort)
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
