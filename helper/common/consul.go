package common

import (
	"time"

	"github.com/hoquangnam45/pharmacy-common-go/microservice/consul"
)

func InitializeConsulClient(advertiseIp string) *consul.Client {
	return consul.SetupEcsConsulClient(func() (map[string]bool, error) {
		return map[string]bool{
			advertiseIp + ":8500": true,
		}, nil
	}, time.Minute*60*24*30)
}
