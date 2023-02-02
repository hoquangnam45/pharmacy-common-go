package ecs

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
)

func GetHostIpAndHostCloudmapPort(ecsMetadataPath string, cloudMapPort int, maxRetryTime time.Duration, interval time.Duration) (string, int, error) {
	if metadata, err := GetEcsMetadata(ecsMetadataPath, maxRetryTime, interval); err == nil {
		return GetHostIpAndHostCloudmapPortFromMetadata(metadata, cloudMapPort)
	} else {
		return "", -1, err
	}
}

func GetHostIpAndHostCloudmapPortFromMetadata(metadata map[string]any, cloudMapPort int) (string, int, error) {
	if status, ok := metadata["MetadataFileStatus"].(string); ok {
		if strings.ToLower(status) == "ready" {
			hostPrivateIp := metadata["HostPrivateIPv4Address"].(string)
			portMappings := metadata["PortMappings"].([]interface{})
			for _, v := range portMappings {
				portMapping, _ := v.(map[string]any)
				containerPort := int(portMapping["ContainerPort"].(float64))
				hostPort := int(portMapping["HostPort"].(float64))
				if containerPort == cloudMapPort {
					return hostPrivateIp, hostPort, nil
				}
			}
		}
	}
	return "", -1, errors.New("metadata file not ready or path is wrong")
}

func GetEcsMetadata(path string, maxRetryTime time.Duration, interval time.Duration) (map[string]any, error) {
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
			return nil, err
		}),
	).RetryUntilSuccess(maxRetryTime, interval).EvalNoCleanup()
}
