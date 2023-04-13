package errorHandler

import (
	"context"
	"time"
)

func (m *MaybeError[T]) DefaultEval(defaultValue T) T {
	val, err := m.Eval()
	if err != nil {
		return defaultValue
	}
	return val
}

func (m *MaybeError[T]) PanicEval() T {
	val, err := m.Eval()
	if err != nil {
		panic(err)
	}
	return val
}

func (m *MaybeError[T]) EvalWithHandler(handler func(error)) (T, error) {
	val, err := m.Eval()
	if err != nil {
		handler(err)
	}
	return val, err
}

func (m *MaybeError[T]) EvalWithHandlerE(handler func(error) error) (T, error) {
	val, err := m.Eval()
	if err != nil {
		innerErr := handler(err)
		if innerErr == nil {
			return val, err
		}
		return val, innerErr
	}
	return val, err
}

func (m *MaybeError[T]) GoEval(retCh chan<- T, errorCh chan<- error) {
	val, err := m.Eval()
	if err != nil {
		errorCh <- err
	} else {
		retCh <- val
	}
	close(errorCh)
	close(retCh)
}

func (m *MaybeError[T]) EvalWithCh() (<-chan T, <-chan error) {
	out := make(chan T, 1)
	errs := make(chan error, 1)
	go m.GoEval(out, errs)
	return out, errs
}

func (m *MaybeError[T]) Eval() (T, error) {
	if m.hasCache {
		return m.result.val, m.result.err
	}
	val, err := m.f()
	if err != nil && m.retry {
		if m.maxRetry > 0 {
			for i := 0; i < m.maxRetry; i++ {
				time.Sleep(m.retryDelay)
				val, err = m.f()
				if err == nil {
					break
				}
			}
		} else if m.maxRetryTime > 0 {
			maxEndTime := time.Now().Add(m.maxRetryTime)
			for time.Now().Before(maxEndTime) {
				time.Sleep(m.retryDelay)
				val, err = m.f()
				if err == nil {
					break
				}
			}
		}
	}
	if m.cache {
		m.result = &result[T]{val, err}
		m.hasCache = true
	}
	return val, err
}

func (m *MaybeError[T]) Error() error {
	_, err := m.Eval()
	return err
}

func (m *MaybeError[T]) EvalWithContext(ctx context.Context) (T, error) {
	dataCh, errCh := m.EvalWithCh()
	var noop T
	select {
	case data := <-dataCh:
		return data, <-errCh
	case <-ctx.Done():
		return noop, ctx.Err()
	}
}
