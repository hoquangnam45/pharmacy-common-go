package common

import (
	"os"
	"time"

	"github.com/hoquangnam45/pharmacy-common-go/helper/ecs"
	"github.com/hoquangnam45/pharmacy-common-go/microservice/consul"
	"github.com/hoquangnam45/pharmacy-common-go/util"
)

var defaultConsulUrlFetcher = func() (map[string]bool, error) {
	ret := map[string]bool{}
	if ecsConsulUrls, ok := os.LookupEnv("ECS_CONSUL_SERVER_URL"); ok {
		urls, err := ecs.ResolveHostModeService(ecsConsulUrls)
		if err != nil {
			return nil, err
		}
		ret = util.MergeMap(urls, ret)
	}
	ret = util.MergeMap(ret, util.ListToSet(util.Tokenize(os.Getenv("CONSUL_SERVER_URLS"), ",")))
	ret["localhost"] = true
	return ret, nil
}

func InitializeConsulClient() *consul.Client {
	return consul.SetupEcsConsulClient(defaultConsulUrlFetcher, time.Minute)
}
