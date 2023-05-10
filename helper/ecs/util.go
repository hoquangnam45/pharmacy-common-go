package ecs

import (
	"context"
	"errors"
	"strings"

	"github.com/hoquangnam45/pharmacy-common-go/util/dns"
	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
	"github.com/hoquangnam45/pharmacy-common-go/util/log"
)

type Ecs struct {
	logger      log.Logger
	dnsResolver *dns.DnsResolver
}

func NewEcs(logger log.Logger, dnsResolver *dns.DnsResolver) *Ecs {
	return &Ecs{logger: logger, dnsResolver: dnsResolver}
}

func (e *Ecs) GetAdvertiseIp(ecsMetadataPath string) (string, error) {
	if ecsMetadataPath != "" {
		containerInfo, err := h.Lift(e.GetContainerInfo)(ecsMetadataPath).Eval()
		if err != nil {
			return "", err
		}
		return containerInfo.HostIp, err
	}
	return "", errors.New("missing ecs metadata path")
}

func (e *Ecs) GetAdvertisePort(ecsMetadataPath string, containerPort int) (int, error) {
	if ecsMetadataPath != "" {
		containerInfo, err := h.Lift(e.GetContainerInfo)(ecsMetadataPath).Eval()
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

func (e *Ecs) ResolveHostModeService(ctx context.Context, srvUrl string) (map[string]bool, error) {
	return h.FlatMap(
		h.FactoryM(func() (map[string]bool, error) {
			return e.dnsResolver.ResolveSrvDns(ctx, srvUrl)
		}),
		h.Lift(func(m map[string]bool) (map[string]bool, error) {
			newMap := map[string]bool{}
			for k := range m {
				parts := strings.Split(k, ":")
				resolvedAddrs, err := e.dnsResolver.ResolveADns(ctx, parts[0])
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
