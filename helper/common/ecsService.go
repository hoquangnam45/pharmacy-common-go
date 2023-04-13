package common

import (
	"os"
	"time"

	"github.com/hoquangnam45/pharmacy-common-go/helper/ecs"
	"github.com/hoquangnam45/pharmacy-common-go/util"
	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
)

func InitializeEcsService(servicePorts ...int) (string, map[int]int, string) {
	advertiseIp := ""
	if ecsMetadataPath, ok := os.LookupEnv("ECS_CONTAINER_METADATA_FILE"); ok {
		advertiseIp = h.Lift(ecs.GetAdvertiseIp)(ecsMetadataPath).RetryUntilSuccess(time.Second*20, time.Second*5).PanicEval()
	} else if bindInterface_, ok := os.LookupEnv("CONSUL_BIND_INTERFACE"); ok {
		advertiseIp = h.Lift(util.FindBindInterfaceAddress)(bindInterface_).PanicEval()
	} else {
		pair := h.FactoryM(util.FindFirstNonLoopBackAddr).PanicEval()
		advertiseIp = pair.First
	}

	advertisePortMap := map[int]int{}
	for _, port := range servicePorts {
		advertisePort := h.FactoryM(func() (int, error) {
			return ecs.GetAdvertisePort(os.Getenv("ECS_CONTAINER_METADATA_FILE"), port)
		}).DefaultEval(port)
		advertisePortMap[port] = advertisePort
	}

	clusterPrefix := os.Getenv("CLUSTER_PREFIX")
	return advertiseIp, advertisePortMap, clusterPrefix
}
