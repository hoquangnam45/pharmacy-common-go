package common

import (
	"time"

	"github.com/hoquangnam45/pharmacy-common-go/microservice/consul"
	"github.com/hoquangnam45/pharmacy-common-go/microservice/lb"
	"github.com/hoquangnam45/pharmacy-common-go/util/log"
)

func InitializeConsulClient(advertiseIp string, logger log.Logger) *consul.Client {
	return consul.SetupEcsConsulClient(lb.NewRandomLB(func() (map[string]bool, error) {
		return map[string]bool{
			advertiseIp + ":8500": true,
		}, nil
	}, time.Minute*60*24*30, logger), logger)
}
