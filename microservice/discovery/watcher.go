package discovery

import (
	"context"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	h "github.com/hoquangnam45/pharmacy-common-go/util/errorHandler"
)

type ConsulWatcher struct {
	serviceName string
	ctx         context.Context
	discovery   Discovery
	stop        bool
	mu          sync.RWMutex
}

func NewWatcher(serviceName string, ctx context.Context, discovery Discovery) registry.Watcher {
	return &ConsulWatcher{
		ctx:         ctx,
		serviceName: serviceName,
		discovery:   discovery,
	}
}

func (w *ConsulWatcher) Next() ([]*registry.ServiceInstance, error) {
	return h.FactoryM(func() ([]*registry.ServiceInstance, error) {
		for !w.IsStop() {
			refreshed, err := w.discovery.RefreshService(w.ctx, w.serviceName)
			if err == nil && refreshed {
				return w.discovery.GetService(w.ctx, w.serviceName)
			}
			time.Sleep(500 * time.Millisecond)
		}
		return w.discovery.GetService(w.ctx, w.serviceName)
	}).EvalWithContext(w.ctx)
}

func (w *ConsulWatcher) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.stop = true
	return nil
}

func (w *ConsulWatcher) IsStop() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.stop
}
