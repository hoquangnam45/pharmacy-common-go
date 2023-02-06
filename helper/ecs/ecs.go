package ecs

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	handler "github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
)

type ContainerInfo struct {
	HostIp       string
	PortMappings map[int]int
}

func GetContainerInfo(ecsMetadataPath string) (*ContainerInfo, error) {
	return handler.FlatMap2(
		handler.Just(ecsMetadataPath),
		handler.Lift(GetEcsMetadata),
		handler.Lift(func(metadata map[string]any) (*ContainerInfo, error) {
			return GetContainerInfoFromMetadata(metadata)
		}),
	).Eval()
}

func GetContainerInfoFromMetadata(metadata map[string]any) (*ContainerInfo, error) {
	if status, ok := metadata["MetadataFileStatus"].(string); ok {
		if strings.ToLower(status) == "ready" {
			containerInfo := &ContainerInfo{}
			containerInfo.HostIp = metadata["HostPrivateIPv4Address"].(string)
			containerInfo.PortMappings = map[int]int{}
			portMappings := metadata["PortMappings"].([]interface{})
			for _, v := range portMappings {
				portMapping, _ := v.(map[string]any)
				containerPort := int(portMapping["ContainerPort"].(float64))
				hostPort := int(portMapping["HostPort"].(float64))
				containerInfo.PortMappings[containerPort] = hostPort
			}
			return containerInfo, nil
		}
	}
	return nil, errors.New("metadata file not ready or path is wrong")
}

func GetEcsMetadata(path string) (map[string]any, error) {
	return handler.FlatMap3(
		handler.Just(path),
		handler.Lift(os.Open),
		handler.Lift(func(file *os.File) ([]byte, error) {
			defer file.Close()
			return io.ReadAll(file)
		}),
		handler.Lift(func(bytes []byte) (map[string]any, error) {
			var ret map[string]any
			err := json.Unmarshal(bytes, &ret)
			return ret, err
		}),
	).Eval()
}


