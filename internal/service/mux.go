package service

import (
	"encoding/json"
)

func (s *Service) Mux(urls []string) (map[string]json.RawMessage, error) {
	result := make(map[string]json.RawMessage, len(urls))

	for _, v := range urls {
		bodyBytes, err := s.hr.Get(v)
		if err != nil {
			return nil, err
		}

		result[v] = bodyBytes
	}

	return result, nil
}
