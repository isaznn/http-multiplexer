package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) MuxerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"allowed only post method", http.StatusMethodNotAllowed)
		return
	}

	_, _ = fmt.Fprintln(w, "Hello World")
}
