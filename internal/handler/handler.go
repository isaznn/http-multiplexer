package handler

import (
	"context"
	"net/http"
)

type Muxer interface {
	Mux(ctx context.Context, urls []string) (map[string]string, error)
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
