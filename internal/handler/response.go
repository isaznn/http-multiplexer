package handler

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

func (h *Handler) errResponse(w http.ResponseWriter, err error, httpStatus int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(apiError{
		Error:        true,
		ErrorMessage: err.Error(),
	})
}
