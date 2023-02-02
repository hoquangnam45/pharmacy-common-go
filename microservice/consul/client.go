package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
)

type ConsulClient struct {
	*api.Client
}

func NewClient(addr string) (*ConsulClient, error) {
	config := api.DefaultConfig()
	config.Address = addr
	return errorHandler.FlatMap2(
		errorHandler.Just(config),
		errorHandler.Lift(api.NewClient),
		func(client *api.Client) *errorHandler.MaybeError[*ConsulClient] {
			return errorHandler.Just(&ConsulClient{
				client,
			})
		},
	).EvalNoCleanup()
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
