package service

type HttpRequester interface {
	Get(url string) ([]byte, error)
}

type Service struct {
	hr HttpRequester
}

func NewService(hr HttpRequester) *Service {
	return &Service{
		hr: hr,
	}
}
