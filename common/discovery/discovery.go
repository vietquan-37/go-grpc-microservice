package discovery

import (
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(instanceID, serviceName, hostPort string) error
	Deregister(instanceID, serviceName string)
	HealthCheck(instanceID string) error
}

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s_%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
