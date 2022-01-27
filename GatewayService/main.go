package main

import (
	"fmt"
	"github.com/heetch/Robert-technical-test/GatewayService/routes"
	"github.com/heetch/Robert-technical-test/service"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"log"
	"net/http"
	"os"
)

func main() {
		tracer.Start()
		defer tracer.Stop()

		Serv, err := service.NewService(os.Getenv("SERVICE_NAME"))

		if err != nil {
			log.Fatalln(err)
		}

     	Serv.Register()

		r := routes.NewRouter(Serv.CSA)

		log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("SERVICE_PORT")), r))
}
