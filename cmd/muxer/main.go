package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/isaznn/http-multiplexer/internal/external"
	"github.com/isaznn/http-multiplexer/internal/handler"
	"github.com/isaznn/http-multiplexer/internal/service"
)

const (
	requestLimit = 100
	concurrentRequestsLimit = 4
	httpClientTimeoutSec = 1
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
		Timeout: httpClientTimeoutSec * time.Second,
	})
	s := service.NewService(concurrentRequestsLimit, ex)
	h := handler.NewHandler(requestLimit, s)
	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%s", srvHost, srvPort),
		Handler: h.InitRouter(),
	}

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Fatalln(err)
		}
		close(idleConnectionsClosed)
	}()
	fmt.Printf("listen on port %s..", srvPort)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalln(err)
	}

	<-idleConnectionsClosed
	fmt.Println("server stopped successfully")
}
