package consul

import (
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hoquangnam45/pharmacy-common-go/microservice/lb"
	"github.com/hoquangnam45/pharmacy-common-go/util"
	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
)

type Client struct {
	client              *api.Client
	config              api.Config
	consulUrlLb         *lb.RandomLB[string]
	serviceUrlLbs       map[string]*lb.RoundRobinLB[string]
	RegisteredAddr      string
	serviceRegistration *api.AgentServiceRegistration
}

func NewClient(config api.Config, consulUrlFetcher lb.ElementFetcher[string], ttl time.Duration) *Client {
	return &Client{
		config:        config,
		consulUrlLb:   lb.NewRandomLB(consulUrlFetcher, ttl),
		serviceUrlLbs: map[string]*lb.RoundRobinLB[string]{},
	}
}

func (c *Client) RefreshServiceUrls(serviceName string) func() (map[string]bool, error) {
	return h.FlatMap(
		h.FactoryM(func() ([]*api.ServiceEntry, error) {
			services, _, err := c.client.Health().Service(serviceName, "", true, nil)
			return services, err
		}),
		h.Lift(func(services []*api.ServiceEntry) (map[string]bool, error) {
			m := map[string]bool{}
			for _, v := range services {
				address := v.Service.Address + ":" + strconv.Itoa(v.Service.Port)
				m[address] = true
			}
			return m, nil
		}),
	).Unwrap()
}

func (c *Client) ConnectToConsul() error {
	loadbalancer := c.consulUrlLb
	err := loadbalancer.CheckRefresh()
	if err != nil {
		return err
	}

	for consulAddr, err := loadbalancer.LoadBalancing(); err == nil; consulAddr, err = loadbalancer.LoadBalancing() {
		config := c.config
		config.Address = consulAddr
		client, err := api.NewClient(&config)
		if err != nil {
			util.SugaredLogger.Infof("Can't initialize consul client at %s. Error %s\n", consulAddr, err.Error())
			util.SugaredLogger.Infof("Delete %s from load balancer", consulAddr)
			loadbalancer.Remove(consulAddr)
			continue
		}
		if err := c.checkConsulConnection(client); err != nil {
			util.SugaredLogger.Infof("Can't register to consul at %s. Error %s\n", consulAddr, err.Error())
			util.SugaredLogger.Infof("Delete %s from load balancer", consulAddr)
			loadbalancer.Remove(consulAddr)
			continue
		} else {
			c.RegisteredAddr = consulAddr
			c.client = client
			return nil
		}
	}
	return errors.New("can't connect to any consul agent")
}

func (c *Client) RegisterService(serviceRegistration *api.AgentServiceRegistration) error {
	if c.serviceRegistration == nil {
		c.serviceRegistration = serviceRegistration
	}
	return c.client.Agent().ServiceRegister(serviceRegistration)
}

func (c *Client) LoadBalancing(serviceName string) (string, error) {
	loadbalancer, ok := c.serviceUrlLbs[serviceName]
	if !ok {
		loadbalancer = lb.NewRoundRobinLB(c.RefreshServiceUrls(serviceName), time.Minute)
		c.serviceUrlLbs[serviceName] = loadbalancer
	}
	return loadbalancer.Get()
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
