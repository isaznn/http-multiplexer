package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/isaznn/http-multiplexer/internal/handler"
)

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

	h := handler.NewHandler()
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", srvHost, srvPort), h.InitRouter())
	if err != nil {
		log.Fatalln(err)
	}
}
