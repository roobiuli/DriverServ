	package main

	import (
		"fmt"
		conapi "github.com/hashicorp/consul/api"
		"log"
		"net/http"
		"os"
	)

	func main() {



		defconf := conapi.DefaultConfig()
		defconf.Address = os.Getenv("CONSUL_HOST") + ":8500"
		con, _ := conapi.NewClient(defconf)

		service, err := ResolveHost("drivers", con)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(service)


		req, _ := http.NewRequest("GET", "http://" + service + "/drivers/42", nil)

		httpc := &http.Client{}

		resp, err :=httpc.Do(req)
		if err != nil{
			fmt.Println(err)
		}

		fmt.Println(resp.StatusCode)



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