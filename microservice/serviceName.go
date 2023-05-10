package microservice

import (
	"errors"
	"fmt"
)

type ServiceNameMapper interface {
	GetServiceName(serviceType string) (string, error)
}

type serviceNameMapper struct {
	services      map[string]bool
	clusterPrefix string
}

func (s *serviceNameMapper) NewServiceNameMapper(clusterPrefix string, services []string) {
	s.services = map[string]bool{}
	for _, service := range services {
		s.services[service] = true
		s.clusterPrefix = clusterPrefix
	}
}

func (s *serviceNameMapper) GetServiceName(serviceType string) (string, error) {
	_, ok := s.services[serviceType]
	if !ok {
		return "", fmt.Errorf("%w %s", errors.New("not supported service type"), serviceType)
	}
	return fmt.Sprintf("%s-%s", s.clusterPrefix, serviceType), nil
}
