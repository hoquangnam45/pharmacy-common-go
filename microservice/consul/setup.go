package consul

import (
	"github.com/hoquangnam45/pharmacy-common-go/microservice/lb"
	"github.com/hoquangnam45/pharmacy-common-go/util/log"
)

func SetupEcsConsulClient(lb lb.LoadBalancer[string], logger log.Logger) *Client {
	return NewClient(lb, logger)
}
