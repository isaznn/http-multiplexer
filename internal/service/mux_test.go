package service

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

const (
	mockGetDelayMs = 30
	mockConcurrentRequestsLimit = 3
)

/*
* Mock Get return JSON with timestamp and received url by default
* and return error if <url> contains path /bad
* with 30 ms delay
*/
type mockExternal struct {}
func (ex *mockExternal) Get(url string) ([]byte, error) {
	// imitation of work
	time.Sleep(mockGetDelayMs * time.Millisecond)

	switch strings.Contains(url, "/bad") {
	case true:
		return nil, fmt.Errorf("request with url %s have status code 400 Bad Request", url)
	default:
		return []byte(fmt.Sprintf(`{"timestamp":%d,"message":"Url: %s"}`, time.Now().Unix(), url)), nil
	}
}

func TestService_Mux(t *testing.T) {
	t.Run("error in one of the urls", func(t *testing.T) {
		// arrange
		urls := []string{"https://example.com", "https://google.com/news", "https://domain.com/with/bad"}
		s := NewService(mockConcurrentRequestsLimit, &mockExternal{})

		// act
		_, err := s.Mux(context.Background(), urls)

		// assert
		if err == nil {
			t.Error("not return an error")
		}
	})

	t.Run("does not processing requests concurrently", func(t *testing.T) {
		// arrange
		timeNow := time.Now()
		urls := []string{"https://example.com", "https://google.com/news", "https://apple.com"}
		s := NewService(mockConcurrentRequestsLimit, &mockExternal{})

		// act
		_, err := s.Mux(context.Background(), urls)

		// assert
		if err != nil {
			t.Errorf("an error has occurred: %v", err)
		}
		if time.Since(timeNow) > time.Duration(len(urls) * mockGetDelayMs + 5) * time.Millisecond {
			t.Error("execution time exceeded")
		}
	})
}
