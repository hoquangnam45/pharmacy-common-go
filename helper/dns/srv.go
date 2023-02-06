package dns

import (
	"log"
	"net"
	"strconv"

	handler "github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
)

func ResolveSrvDns(link string) (map[string]bool, error) {
	return handler.FlatMap2(
		handler.Just(link),
		handler.Lift(func(host string) ([]*net.SRV, error) {
			log.Printf("Start lookingup host %s", host)
			_, addrs, err := net.LookupSRV("", "", host)
			if err != nil {
				return nil, err
			}
			return addrs, nil
		}),
		handler.Lift(func(addrs []*net.SRV) (map[string]bool, error) {
			resolvedAddrs := map[string]bool{}
			log.Printf("Found %d records: ", len(addrs))
			for _, v := range addrs {
				resolvedAddr := v.Target + ":" + strconv.Itoa(int(v.Port))
				log.Print(resolvedAddr)
				resolvedAddrs[resolvedAddr] = true
			}
			return resolvedAddrs, nil
		}),
	).Eval()
}
