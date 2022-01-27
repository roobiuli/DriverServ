package service

import (
	"fmt"
	conapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
	"github.com/heetch/Robert-technical-test/GatewayService/routes"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)
// Service struct used to create, register new services
type Service struct {
	Name string
	CSA *conapi.Client
	Router *muxtrace.Router
}


func NewService(n string) (*Service, error) {
	conf := conapi.DefaultConfig()
	conf.Address = os.Getenv("CONSUL_HOST") + ":8500"
	consulapi, err := conapi.NewClient(conf)

	if err != nil {
		return nil, err
	}

	return &Service{
		Name:   n,
		CSA:    consulapi,
		Router: routes.NewRouter(consulapi),
	}, nil
}


func (s *Service) Register() {
	// Registering service with Consul

	registration := new(conapi.AgentServiceRegistration)
	registration.Name = s.Name
	address := hostname()
	registration.Address = address
	port, err := strconv.Atoi(port()[1:len(port())])
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = port
	registration.Check = new(conapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck",
		address, port)
	registration.Check.Interval = "30s"
	registration.Check.Timeout = "3s"

	s.CSA.Agent().ServiceRegister(registration)

}

func (s *Service) DeRegister(id string) error  {
	return s.CSA.Agent().ServiceDeregister(id)
}



func (s *Service) ServiceStart(port string) {
	// Not used, developed for Consul Connect Service[Discovery, communication] communication via TLS
	svc, _ := connect.NewService(s.Name, s.CSA)
	defer svc.Close()

	server := &http.Server{
		Addr: port,
		Handler: s.Router,
		TLSConfig: svc.ServerTLSConfig(),
	}

	server.ListenAndServeTLS("", "")
}




func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("I'm a healthy service !"))
	return
}


func port() string {
	p := os.Getenv("SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8100"
	}
	return fmt.Sprintf(":%s", p)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}