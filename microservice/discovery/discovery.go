package discovery

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	"github.com/hoquangnam45/pharmacy-common-go/util"
	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
)

type Discovery interface {
	RefreshService(ctx context.Context, serviceName string) (bool, error)
	registry.Discovery
}

type ConsulDiscovery struct {
	client   *api.Client
	services *util.SyncMap[string, []*registry.ServiceInstance]
}

func NewDiscovery(c *api.Client) Discovery {
	return &ConsulDiscovery{
		client:   c,
		services: util.NewSyncMap[string, []*registry.ServiceInstance](),
	}
}

func (d *ConsulDiscovery) RefreshService(ctx context.Context, serviceName string) (bool, error) {
	newServices, err := h.FactoryM(func() ([]*registry.ServiceInstance, error) {
		services, _, err := d.client.Health().Service(serviceName, "", true, nil)
		if err != nil {
			return nil, err
		}
		serviceInstances := []*registry.ServiceInstance{}
		for _, service := range services {
			tags := map[string]string{}
			endpointMap := map[string]bool{}
			for _, tag := range service.Service.Tags {
				tokens := strings.SplitN(tag, "=", 2)
				k := strings.Trim(tokens[0], " ")
				v := strings.Trim(tokens[1], " ")
				tags[k] = v
			}
			if endp, err := util.BuildEndpoint(service.Service.Address, service.Service.Port); err != nil {
				return nil, err
			} else {
				endpointMap[endp] = true
			}
			for _, v := range service.Service.TaggedAddresses {
				port := 0
				address := ""
				scheme := ""
				if url, err := url.Parse(v.Address); err != nil {
					return nil, err
				} else {
					port = v.Port
					address = url.Host
					scheme = url.Scheme
				}
				if endpoint, err := util.BuildEndpointScheme(scheme, address, port); err != nil {
					return nil, err
				} else {
					endpointMap[endpoint] = true
				}
			}
			endpoints := []string{}
			for endp := range endpointMap {
				endpoints = append(endpoints, endp)
			}
			instance := &registry.ServiceInstance{
				ID:        service.Service.ID,
				Name:      service.Service.Service,
				Version:   tags["version"],
				Metadata:  service.Service.Meta,
				Endpoints: endpoints,
			}
			serviceInstances = append(serviceInstances, instance)
		}
		return serviceInstances, nil
	}).EvalWithContext(ctx)
	if err != nil {
		return false, err
	}
	if d.checkServiceChange(newServices) {
		d.services.Set(serviceName, newServices)
		return true, nil
	}
	return false, nil
}

func (d *ConsulDiscovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return NewWatcher(serviceName, ctx, d), nil
}

func (d *ConsulDiscovery) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return d.services.GetOrSet(serviceName, []*registry.ServiceInstance{}), nil
}

func (d *ConsulDiscovery) checkServiceChange(services []*registry.ServiceInstance) bool {
	if d.services.Len() != len(services) {
		return true
	}
	diff := false
	for _, service := range services {
		if diff = !service.Equal(d.services.Get(service.ID)); diff {
			break
		}
	}
	return diff
}
