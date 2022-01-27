package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/heetch/Robert-technical-test/midlewares"
	"github.com/heetch/Robert-technical-test/models"
	"github.com/heetch/Robert-technical-test/service"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var drvdata = make(map[string]string)


func ReadCSV(file string)  {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Unable to open CSV File", file, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal("Unable to parse CSV File", file, err)
	}
	 for _, dr := range records {
		 drvdata[dr[0]] = dr[1]
	 }
}


/// Handler Func to take care of the drivers.CSV

func DRVHandlerFunc(w http.ResponseWriter, r *http.Request)  {
	reqvars := mux.Vars(r)

	if driver, ok := drvdata[reqvars["id"]] ; ok {
		driver := &models.Driver{Id: reqvars["id"], Name: driver}
		jdata, err := json.Marshal(driver)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jdata)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	return
}


func main()  {
	// Tracing with DataDog
	tracer.Start()
	defer tracer.Stop()


	Serv, err := service.NewService(os.Getenv("SERVICE_NAME"))

	if err != nil {
		log.Fatalln(err)
	}

	Serv.Register()

	// MUX with Tracing support from DATADOG
	r := muxtrace.NewRouter()

	// Service type satisfies the handler because of the ServeHTTP method so healthcheck to self
	r.Handle("/healthcheck", Serv)

	r.HandleFunc("/drivers/{id:[0-9]+}", midlewares.LoginMiddleWare(DRVHandlerFunc))

	ReadCSV("./drivers.csv")

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("SERVICE_PORT")), r))


}
