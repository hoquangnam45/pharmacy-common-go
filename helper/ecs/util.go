package ecs

import (
	"errors"
	"strings"

	"github.com/hoquangnam45/pharmacy-common-go/helper/dns"
	h "github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
)

func GetAdvertiseIp(ecsMetadataPath string) (string, error) {
	if ecsMetadataPath != "" {
		containerInfo, err := h.Lift(GetContainerInfo)(ecsMetadataPath).Eval()
		if err != nil {
			return "", err
		}
		return containerInfo.HostIp, err
	}
	return "", errors.New("missing ecs metadata path")
}

func GetAdvertisePort(ecsMetadataPath string, containerPort int) (int, error) {
	if ecsMetadataPath != "" {
		containerInfo, err := h.Lift(GetContainerInfo)(ecsMetadataPath).Eval()
		if err != nil {
			return 0, err
		}
		if hostPort, ok := containerInfo.PortMappings[containerPort]; ok {
			return hostPort, nil
		}
		return 0, errors.New("is using host network mode or wrong container port used")
	}
	return 0, errors.New("missing ecs metadata path")
}

func ResolveHostModeService(srvUrl string) (map[string]bool, error) {
	return h.FlatMap(
		h.Lift(dns.ResolveSrvDns)(srvUrl),
		h.Lift(func(m map[string]bool) (map[string]bool, error) {
			newMap := map[string]bool{}
			for k := range m {
				parts := strings.Split(k, ":")
				resolvedAddrs, err := dns.ResolveADns(parts[0])
				if err != nil {
					return nil, err
				}
				for k2 := range resolvedAddrs {
					newMap[k2] = true
				}
			}
			return newMap, nil
		})).Eval()
}
