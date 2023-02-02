package ecs

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
)

type ContainerInfo struct {
	HostIp       string
	PortMappings map[int]int
}

func GetContainerInfo(ecsMetadataPath string) (*ContainerInfo, error) {
	return errorHandler.FlatMap2(
		errorHandler.Just(ecsMetadataPath),
		errorHandler.Lift(GetEcsMetadata),
		errorHandler.Lift(func(metadata map[string]any) (*ContainerInfo, error) {
			return GetContainerInfoFromMetadata(metadata)
		}),
	).EvalNoCleanup()
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
	return errorHandler.FlatMap3(
		errorHandler.Just(path),
		errorHandler.Lift(os.Open),
		errorHandler.Lift(func(file *os.File) ([]byte, error) {
			defer file.Close()
			return io.ReadAll(file)
		}),
		errorHandler.Lift(func(bytes []byte) (map[string]any, error) {
			var ret map[string]any
			err := json.Unmarshal(bytes, &ret)
			return ret, err
		}),
	).EvalNoCleanup()
}
