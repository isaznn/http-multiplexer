package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type APIError struct {
	Error        bool   `json:"error,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func sendHttpErr(w http.ResponseWriter, err error, httpStatus int) {
	r, _ := json.Marshal(APIError{
		Error:        true,
		ErrorMessage: err.Error(),
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, _ = w.Write(r)
}

func MuxerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendHttpErr(w, fmt.Errorf("allowed only post method"), http.StatusMethodNotAllowed)
		return
	}

	_, _ = fmt.Fprintln(w, "Hello World")
}

func main()  {
	srvHost := "0.0.0.0"
	srvPort := "8080"
	pathHost, pHostExists := os.LookupEnv("SRV_HOST")
	if pHostExists {
		srvHost = pathHost
	}
	pathPort, pPortExists := os.LookupEnv("SRV_PORT")
	if pPortExists {
		srvPort = pathPort
	}

	http.HandleFunc("/muxer", MuxerHandler)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", srvHost, srvPort), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
