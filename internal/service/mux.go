package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	contextErrorText = "canceled with context"
	requestErrorText = "request ended with an error"
	IncDelta = 1
)

type safeMap struct {
	M map[string]string
	sync.Mutex
}

func newSafeMap(l int) *safeMap {
	return &safeMap{
		M: make(map[string]string, l),
	}
}

func (s *safeMap) AllEntities() map[string]string {
	s.Lock()
	defer s.Unlock()
	return s.M
}

func (s *safeMap) Store(key string, value string) {
	s.Lock()
	defer s.Unlock()
	s.M[key] = value
}

func (s *Service) chunks(urls []string) [][]string {
	var dividedUrls [][]string
	chunkSize := (len(urls) + s.concurrentRequestsLimit - 1) / s.concurrentRequestsLimit

	for i := 0; i < len(urls); i += chunkSize {
		end := i + chunkSize
		if end > len(urls) {
			end = len(urls)
		}

		dividedUrls = append(dividedUrls, urls[i:end])
	}

	return dividedUrls
}

func (s *Service) Mux(ctx context.Context, urls []string) (map[string]string, error) {
	var (
		m = newSafeMap(len(urls))
		chunks = s.chunks(urls)
		errCtx, cancel = context.WithCancel(context.Background())
		errCh = make(chan struct{})
		errCounter int32
		ctxCancelCounter int32
		wg sync.WaitGroup
	)

	// if received error - cancel context
	go func() {
		<-errCh
		cancel()
	}()

	// push chunks to goroutines
	for _, chunk := range chunks {
		wg.Add(1)
		go func(urls []string) {
			defer wg.Done()
			for _, url := range urls {
				select {
				case <-ctx.Done():
					atomic.AddInt32(&ctxCancelCounter, IncDelta)
					return
				case <-errCtx.Done():
					return
				default:
					bodyBytes, err := s.Get(url)
					if err != nil {
						atomic.AddInt32(&errCounter, IncDelta)
						errCh <- struct{}{}
						return
					}
					m.Store(url, string(bodyBytes))
				}
			}
		}(chunk)
	}
	wg.Wait()

	if ctxCancelCounter > 0 {
		return nil, fmt.Errorf(contextErrorText)
	}
	if errCounter > 0 {
		return nil, fmt.Errorf(requestErrorText)
	}

	return m.AllEntities(), nil

}
