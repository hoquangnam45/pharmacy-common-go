package consul

import (
	"github.com/hashicorp/consul/api"
	handler "github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
)

type ConsulClient struct {
	*api.Client
}

func NewClient(addr string) (*ConsulClient, error) {
	config := api.DefaultConfig()
	config.Address = addr
	return handler.FlatMap2(
		handler.Just(config),
		handler.Lift(api.NewClient),
		func(client *api.Client) *handler.MaybeError[*ConsulClient] {
			return handler.Just(&ConsulClient{
				client,
			})
		},
	).Eval()
}

func (c *ConsulClient) Register(id string, healthCheckPath string) error {
	return c.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:   id,
		Name: id + "_ms",
		Check: &api.AgentServiceCheck{
			Interval: "30s",
			Timeout:  "60s",
			HTTP:     healthCheckPath,
		},
	})
}
