package util

import (
	"fmt"
	"net/url"

	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
)

func BuildEndpoint(addr string, port int) (string, error) {
	return h.FlatMap(
		h.Lift(url.Parse)(fmt.Sprintf("%s:%d", addr, port)),
		h.LiftJ(func(url *url.URL) string {
			return url.String()
		}),
	).Eval()
}

func BuildEndpointScheme(scheme string, addr string, port int) (string, error) {
	return h.FlatMap(
		h.Lift(url.Parse)(fmt.Sprintf("%s://%s:%d", scheme, addr, port)),
		h.LiftJ(func(url *url.URL) string {
			return url.String()
		}),
	).Eval()
}
