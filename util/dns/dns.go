package dns

import "github.com/hoquangnam45/pharmacy-common-go/util/log"

type DnsResolver struct {
	logger log.Logger
}

func NewDnsResolver(logger log.Logger) *DnsResolver {
	return &DnsResolver{
		logger: logger,
	}
}
