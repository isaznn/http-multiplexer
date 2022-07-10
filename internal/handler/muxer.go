package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type muxerHandlerRequest struct {
	Urls []string `json:"urls"`
}

type urlWithResponse struct {
	Url      string          `json:"url"`
	Response json.RawMessage `json:"response"`
}

type muxerHandlerResponse struct {
	Result []urlWithResponse `json:"result"`
}

func (h *Handler) MuxerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"allowed only post method", http.StatusMethodNotAllowed)
		return
	}

	_, _ = fmt.Fprintln(w, "Hello World")
}
