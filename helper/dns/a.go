package dns

import (
	"fmt"
	"net"

	handler "github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
	"github.com/hoquangnam45/pharmacy-common-go/util"
	"go.uber.org/zap"
)

func ResolveADns(link string) (map[string]bool, error) {
	return handler.FlatMap2(
		handler.Just(link),
		handler.Lift(func(host string) ([]net.IP, error) {
			util.SugaredLogger.Infof("Start lookingup host %s", host)
			addrs, err := net.LookupIP(host)
			if err != nil {
				return nil, err
			}
			return addrs, nil
		}),
		handler.Lift(func(addrs []net.IP) (map[string]bool, error) {
			resolvedAddrs := map[string]bool{}
			for _, v := range addrs {
				resolvedAddr := v.String()
				resolvedAddrs[resolvedAddr] = true
			}
			util.Logger.Info(fmt.Sprintf("Found %d records", len(addrs)),
				zap.Strings("records", util.SetToList(resolvedAddrs)),
			)
			return resolvedAddrs, nil
		}),
	).Eval()
}
