package midlewares

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

func init()  {

	f, er := os.OpenFile(os.Getenv("LOG_FILE"), os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if er != nil {
		log.Fatalln(er)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

type Repose struct {
	status int
	size int
}

type LogRW struct {
	http.ResponseWriter
	Response *Repose
}

func (l *LogRW) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.Response.size += size
	return size, err
}



func (l *LogRW) WriteHeader(statusc int) {
	l.ResponseWriter.WriteHeader(statusc)
	l.Response.status = statusc
}



func LoginMiddleWare(s http.HandlerFunc)  http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start :=  time.Now()

		response := &Repose{
			status: 0,
			size:   0,
		}
		lrw := &LogRW{
			ResponseWriter: w,
			Response:       response,
		}
		s.ServeHTTP(lrw, r)

		duration := time.Since(start)

		log.WithFields(log.Fields{
			"Duration": duration,
			"Service": "DummyService",
			"Status": response.status,
			"uri": r.RequestURI,
			"Method": r.Method,
			"Size": response.size,
		}).Info("HTTP LOG")

	})
}
