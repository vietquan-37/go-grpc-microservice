package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(instanceID, serviceName, hostPort string) error
	Deregister(instanceID, serviceName string)
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceID string) error
}

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s_%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
