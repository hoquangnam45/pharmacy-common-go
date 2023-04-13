package ecs

import (
	"errors"
	"os"
	"strings"

	"github.com/hoquangnam45/pharmacy-common-go/util"
	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
)

type ContainerInfo struct {
	HostIp       string
	PortMappings map[int]int
}

func GetContainerInfo(ecsMetadataPath string) (*ContainerInfo, error) {
	return h.FlatMap(
		h.Lift(GetEcsMetadata)(ecsMetadataPath),
		h.Lift(GetContainerInfoFromMetadata),
	).Eval()
}

func GetContainerInfoFromMetadata(metadata map[string]any) (*ContainerInfo, error) {
	if status, ok := metadata["MetadataFileStatus"].(string); ok {
		if strings.ToLower(status) == "ready" {
			containerInfo := &ContainerInfo{}
			containerInfo.HostIp = metadata["HostPrivateIPv4Address"].(string)
			containerInfo.PortMappings = map[int]int{}
			// Ecs service host network mode check
			if metadata["PortMappings"] != nil {
				portMappings := metadata["PortMappings"].([]interface{})
				for _, v := range portMappings {
					portMapping, _ := v.(map[string]any)
					containerPort := int(portMapping["ContainerPort"].(float64))
					hostPort := int(portMapping["HostPort"].(float64))
					containerInfo.PortMappings[containerPort] = hostPort
				}
			}
			return containerInfo, nil
		}
	}
	return nil, errors.New("metadata file not ready or path is wrong")
}

func GetEcsMetadata(path string) (map[string]any, error) {
	return h.FlatMap2(
		h.Lift(os.Open)(path),
		h.Lift(util.ReadAllThenClose[*os.File]),
		h.Lift(util.UnmarshalJsonDeref(&map[string]any{})),
	).Eval()
}
