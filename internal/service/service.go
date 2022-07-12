package service

type HttpRequester interface {
	Get(url string) ([]byte, error)
}

type Service struct {
	concurrentRequestsLimit int
	HttpRequester
}

func NewService(concurrentRequestsLimit int, hr HttpRequester) *Service {
	return &Service{
		concurrentRequestsLimit: concurrentRequestsLimit,
		HttpRequester:           hr,
	}
}
