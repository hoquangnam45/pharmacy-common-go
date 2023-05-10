package dns

import (
	"context"
	"fmt"
	"net"

	"github.com/hoquangnam45/pharmacy-common-go/util"
	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
	"go.uber.org/zap"
)

func (r *DnsResolver) ResolveADns(ctx context.Context, link string) (map[string]bool, error) {
	return h.FlatMap2(
		h.Just(link),
		h.Lift(func(host string) ([]net.IP, error) {
			r.logger.Info("Start lookingup host %s", host)
			addrs, err := net.LookupIP(host)
			if err != nil {
				return nil, err
			}
			return addrs, nil
		}),
		h.Lift(func(addrs []net.IP) (map[string]bool, error) {
			resolvedAddrs := map[string]bool{}
			for _, v := range addrs {
				resolvedAddr := v.String()
				resolvedAddrs[resolvedAddr] = true
			}
			r.logger.Info(fmt.Sprintf("Found %d records", len(addrs)),
				zap.Strings("records", util.SetToList(resolvedAddrs)),
			)
			return resolvedAddrs, nil
		}),
	).EvalWithContext(ctx)
}
