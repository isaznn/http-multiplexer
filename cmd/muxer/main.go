package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/isaznn/http-multiplexer/internal/external"
	"github.com/isaznn/http-multiplexer/internal/handler"
	"github.com/isaznn/http-multiplexer/internal/service"
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

	ex := external.NewExternal(&http.Client{
		Timeout: 1 * time.Second,
	})
	s := service.NewService(ex)
	h := handler.NewHandler(s)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", srvHost, srvPort), h.InitRouter())
	if err != nil {
		log.Fatalln(err)
	}
}
