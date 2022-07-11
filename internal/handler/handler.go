package handler

import (
	"encoding/json"
	"net/http"
)

type Muxer interface {
	Mux(urls []string) (map[string]json.RawMessage, error)
}

type Handler struct {
	requestLimit int32
	requestCounter int32
	Muxer
}

func NewHandler(requestLimit int32, m Muxer) *Handler {
	return &Handler{
		requestLimit:   requestLimit,
		requestCounter: 0,
		Muxer:          m,
	}
}

func (h *Handler) InitRouter() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/muxer", h.muxHandler)
	return r
}
