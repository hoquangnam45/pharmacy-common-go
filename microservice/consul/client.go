package consul

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"
	"github.com/hoquangnam45/pharmacy-common-go/microservice"
	"github.com/hoquangnam45/pharmacy-common-go/microservice/discovery"
	"github.com/hoquangnam45/pharmacy-common-go/microservice/lb"
	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
	"github.com/hoquangnam45/pharmacy-common-go/util/log"
)

type Client struct {
	client      *api.Client
	consulUrlLb lb.LoadBalancer[string]
	discovery.Discovery
	*KVClient
	logger            log.Logger
	clients           map[string]*http.Client
	serviceNameMapper microservice.ServiceNameMapper
}

func NewClient(lb lb.LoadBalancer[string], logger log.Logger) *Client {
	return &Client{
		consulUrlLb: lb,
		logger:      logger,
	}
}

func WrapBaseClient(c *api.Client) *Client {
	return &Client{
		client:    c,
		Discovery: discovery.NewDiscovery(c),
		KVClient:  NewKvClient(c),
	}
}

func (c *Client) ConnectToConsul(config *api.Config) error {
	loadbalancer := c.consulUrlLb
	err := loadbalancer.CheckRefresh()
	if err != nil {
		return err
	}

	for consulAddr, err := loadbalancer.LoadBalancing(); err == nil; consulAddr, err = loadbalancer.LoadBalancing() {
		config := *config
		config.Address = consulAddr
		client, err := api.NewClient(&config)
		if err != nil {
			c.logger.Info("Can't initialize consul client at %s. Error %s\n", consulAddr, err.Error())
			c.logger.Info("Delete %s from load balancer", consulAddr)
			loadbalancer.Remove(consulAddr)
			continue
		}
		if err := c.checkConsulConnection(client); err != nil {
			c.logger.Info("Can't register to consul at %s. Error %s\n", consulAddr, err.Error())
			c.logger.Info("Delete %s from load balancer", consulAddr)
			loadbalancer.Remove(consulAddr)
			continue
		} else {
			c.client = client
			c.KVClient = NewKvClient(client)
			c.Discovery = discovery.NewDiscovery(c.client)
			return nil
		}
	}
	return errors.New("can't connect to any consul agent")
}

func (c *Client) RegisterService(serviceRegistration *api.AgentServiceRegistration) error {
	return c.client.Agent().ServiceRegister(serviceRegistration)
}

func (c *Client) CheckConsulConnection() error {
	return c.checkConsulConnection(c.client)
}

func (c *Client) checkConsulConnection(client *api.Client) error {
	if _, err := client.Status().Peers(); err != nil {
		return err
	}
	return nil
}

func (c *Client) CallService(ctx context.Context, serviceType string, method string, path string, args interface{}, reply interface{}) error {
	return h.FlatMap(
		h.Lift(c.serviceNameMapper.GetServiceName)(serviceType),
		h.LiftE(func(serviceName string) error {
			client, ok := c.clients[serviceName]
			if !ok {
				clientI, err := http.NewClient(context.Background(),
					http.WithEndpoint(fmt.Sprintf("discovery:///%s", serviceName)),
					http.WithDiscovery(c.Discovery))
				if err != nil {
					return err
				}
				client = clientI
				c.clients[serviceName] = client
			}
			return client.Invoke(ctx, method, path, args, reply)
		}),
	).Error()
}

// func (c *Client) RefreshServiceUrls(serviceName string) func() (map[string]bool, error) {
// 	return h.FlatMap(
// 		h.FactoryM(func() ([]*api.ServiceEntry, error) {
// 			services, _, err := c.client.Health().Service(serviceName, "", true, nil)
// 			return services, err
// 		}),
// 		h.Lift(func(services []*api.ServiceEntry) (map[string]bool, error) {
// 			m := map[string]bool{}
// 			for _, v := range services {
// 				address := v.Service.Address + ":" + strconv.Itoa(v.Service.Port)
// 				m[address] = true
// 			}
// 			return m, nil
// 		}),
// 	).Unwrap()
// }

// func (c *Client) LoadBalancing(serviceName string) (string, error) {
// 	loadbalancer, ok := c.serviceUrlLbs[serviceName]
// 	if !ok {
// 		loadbalancer = lb.NewRoundRobinLB(c.RefreshServiceUrls(serviceName), time.Minute)
// 		c.serviceUrlLbs[serviceName] = loadbalancer
// 	}
// 	return loadbalancer.Get()
// }

// func (c *Client) GetClient(serviceName string) (error) {
// 	client, exist := c.serviceClients.Clone().ExistGet(serviceName)
// 	if !exist {
// 		// Create a route Filter: filter instances with version number "2.0.0"
// 		filters := filter.Version("2.0.0")
// 		filter.Version
// 		// Create P2C load balancing algorithm Selector, and inject routing Filter
// 		client, err := http.NewClient(
// 			context.Background(),
// 			http.WithNodeFilter(filters),
// 			http.WithEndpoint(fmt.Sprintf("discovery:///%s", serviceName)),
// 			http.WithDiscovery(c.Discovery),
// 		)
// 	}
// }
