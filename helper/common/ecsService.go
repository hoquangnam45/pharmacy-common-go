package common

import (
	"os"
	"time"

	"github.com/hoquangnam45/pharmacy-common-go/helper/ecs"
	h "github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
	"github.com/hoquangnam45/pharmacy-common-go/helper/net"
)

func InitializeEcsService(servicePort int) (string, int, string) {
	advertiseIp := ""
	if ecsMetadataPath, ok := os.LookupEnv("ECS_CONTAINER_METADATA_FILE"); ok {
		advertiseIp = h.FactoryM(func() (string, error) {
			return ecs.GetAdvertiseIp(ecsMetadataPath)
		}).RetryUntilSuccess(time.Second*20, time.Second*5).PanicEval()
	} else if bindInterface_, ok := os.LookupEnv("CONSUL_BIND_INTERFACE"); ok {
		advertiseIp = h.Lift(net.FindBindInterfaceAddress)(bindInterface_).PanicEval()
	} else {
		pair := h.FactoryM(net.FindFirstNonLoopBackAddr).PanicEval()
		advertiseIp = pair.First
	}

	advertisePort := h.FactoryM(func() (int, error) {
		return ecs.GetAdvertisePort(os.Getenv("ECS_CONTAINER_METADATA_FILE"), servicePort)
	}).DefaultEval(servicePort)

	clusterPrefix := os.Getenv("CLUSTER_PREFIX")
	return advertiseIp, advertisePort, clusterPrefix
}
