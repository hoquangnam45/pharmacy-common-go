package consul

import (
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
	"github.com/hoquangnam45/pharmacy-common-go/microservice/lb"
)

type ConsulClient struct {
	consulUrlsRefresher func() (map[string]bool, error)
	client              *api.Client
	config              api.Config
	consulUrlLb         *lb.RandomLB[string]
	serviceUrlLbs       map[string]*lb.RoundRobinLB[string]
}

func NewClient(config api.Config, consulUrlsRefresher func() (map[string]bool, error)) *ConsulClient {
	return &ConsulClient{
		consulUrlsRefresher: consulUrlsRefresher,
		config:              config,
		consulUrlLb:         lb.NewRandomLB[string](),
		serviceUrlLbs:       map[string]*lb.RoundRobinLB[string]{},
	}
}

func (c *ConsulClient) RefreshServiceUrls(serviceName string) (map[string]bool, error) {
	return errorHandler.FlatMap(
		errorHandler.FactoryM(func() ([]*api.CatalogService, error) {
			services, _, err := c.client.Catalog().Service(serviceName, "", nil)
			return services, err
		}),
		errorHandler.Lift(func(services []*api.CatalogService) (map[string]bool, error) {
			m := map[string]bool{}
			for _, v := range services {
				address := v.ServiceAddress + ":" + strconv.Itoa(v.ServicePort)
				m[address] = true
			}
			return m, nil
		}),
	).Eval()
}

func (c *ConsulClient) Register(serviceRegistration *api.AgentServiceRegistration) error {
	loadbalancer := c.consulUrlLb
	err := loadbalancer.Check()
	if err == lb.ErrorEmptyList || err == lb.ErrorNeedRefresh {
		_, err = errorHandler.FlatMap(
			errorHandler.FactoryM(c.consulUrlsRefresher),
			errorHandler.PeekE(func(consulUrls map[string]bool) error {
				loadbalancer.RefreshList(consulUrls, time.Minute)
				return loadbalancer.Check()
			})).Eval()
	}
	if err != nil {
		return err
	}

	for consulAddr, err := loadbalancer.Get(); err == nil; consulAddr, err = loadbalancer.Get() {
		config := c.config
		config.Address = consulAddr
		client, err_ := errorHandler.FlatMap(
			errorHandler.Lift(api.NewClient)(&config),
			errorHandler.PeekE(func(client *api.Client) error {
				return client.Agent().ServiceRegister(serviceRegistration)
			})).Eval()
		if err_ != nil {
			loadbalancer.Remove(config.Address)
			continue
		}
		c.client = client
		return nil
	}
	return errors.New("can't register with any consul urls")
}

func (c *ConsulClient) LoadBalancing(serviceName string) (string, error) {
	loadbalancer, ok := c.serviceUrlLbs[serviceName]
	if !ok {
		c.serviceUrlLbs[serviceName] = lb.NewRoundRobinLB[string]()
		loadbalancer = c.serviceUrlLbs[serviceName]
	}
	err := loadbalancer.Check()
	if err == lb.ErrorEmptyList || err == lb.ErrorNeedRefresh {
		_, err = errorHandler.FlatMap(
			errorHandler.Lift(c.RefreshServiceUrls)(serviceName),
			errorHandler.PeekE(func(serviceUrls map[string]bool) error {
				loadbalancer.RefreshList(serviceUrls, time.Minute)
				return loadbalancer.Check()
			})).Eval()
	}
	if err != nil {
		return "", err
	}
	return loadbalancer.Get()
}
