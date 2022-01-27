package routes

import (
	"fmt"
	conapi "github.com/hashicorp/consul/api"
	"github.com/heetch/Robert-technical-test/midlewares"
	"io/ioutil"
	"log"
	"os"

	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	//"io/ioutil"
	"net/http"
	"strings"
)
// Consul Host for service autodiscovery

var ConsAPI *conapi.Client

// Rute Struct holding all info about a route

func init() {
	defConf := conapi.DefaultConfig()
	defConf.Address = os.Getenv("CONSUl_HOST") + ":8500"
	Capi, err := conapi.NewClient(defConf)
	if err != nil {
		log.Println(err)
	}

	ConsAPI = Capi

}
type Route struct {
	Name string
	Path string
	Handler http.HandlerFunc
	Method string
}

// Intended to be looped in GO MUX to automatically spawn the routes at RUNTIME

type Routes []*Route


// one can append more/new routes for future micro services

var routes = Routes{

	&Route{
		Name:    "drivers",
		Path:    "/drivers/{id:[0-9]+}",
		Handler: drivers,
		Method:  "GET",
		},
	&Route{
		Name:    "HealthCheck",
		Path:    "/healthcheck",
		Handler: healthcheck,
		Method:  "GET",
	},
}

func NewRouter(client *conapi.Client) *muxtrace.Router {
	ConsAPI = client
	r := muxtrace.NewRouter().StrictSlash(true)
	for _, route := range routes {
		r.Methods(route.Method).HandlerFunc(midlewares.LoginMiddleWare(route.Handler)).Path(route.Path).Name(route.Name)
	}
	return r
}



func drivers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Allowed only GET Method"))
	}

	// Ugly way of doing Service discovery
	fmt.Println(r.URL.Path)
	srvname := strings.Split(r.URL.Path, "/")

	service, err := ResolveHost(srvname[1],ConsAPI)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v%v", service,  r.URL.Path ), nil)

	req.Header = make(http.Header)
	req.Header.Set("Host", r.Host)
	req.Header.Set("X-Forwarded-For", r.RemoteAddr)

	for h, val := range r.Header {
		req.Header[h] = val
	}

	hcl := http.Client{}

	resp, err := hcl.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()


	respbody, err := ioutil.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(respbody)
	return

}




func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service IS OK"))

}


func ResolveHost(sname string, client *conapi.Client ) (string, error) {

	serv, _, err := client.Health().Service(sname, "", true, nil)
	if err != nil {
		return "", err
	}
	srv := serv[0]
	service := fmt.Sprintf("%v:%v", srv.Service.Address, srv.Service.Port)

	return service, nil
}