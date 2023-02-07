package dns

import (
	"log"
	"net"

	handler "github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
)

func ResolveADns(link string) (map[string]bool, error) {
	return handler.FlatMap2(
		handler.Just(link),
		handler.Lift(func(host string) ([]net.IP, error) {
			log.Printf("Start lookingup host %s", host)
			addrs, err := net.LookupIP(host)
			if err != nil {
				return nil, err
			}
			return addrs, nil
		}),
		handler.Lift(func(addrs []net.IP) (map[string]bool, error) {
			resolvedAddrs := map[string]bool{}
			log.Printf("Found %d records: ", len(addrs))
			for _, v := range addrs {
				resolvedAddr := v.String()
				log.Print(resolvedAddr)
				resolvedAddrs[resolvedAddr] = true
			}
			return resolvedAddrs, nil
		}),
	).Eval()
}
