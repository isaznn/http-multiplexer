package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func MuxerHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Hello World")
}

func main()  {
	srvPort := "8080"
	pathPort, exists := os.LookupEnv("SRV_PORT")
	if exists {
		srvPort = pathPort
	}

	http.HandleFunc("/muxer", MuxerHandler)

	err := http.ListenAndServe(":" + srvPort, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
