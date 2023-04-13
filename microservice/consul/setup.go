package consul

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hoquangnam45/pharmacy-common-go/microservice/lb"
)

func SetupEcsConsulClient(consulUrlFetcher lb.ElementFetcher[string], ttl time.Duration) *Client {
	config := api.DefaultConfig()
	return NewClient(*config, consulUrlFetcher, ttl)
}
