package main

import (
	"common"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	disco "github.com/slink-go/disco-go"
	"github.com/slink-go/logger"
	"net/http"
	"strings"
	"syscall"
	"time"
)

const serviceHost = "localhost"
const serviceProto = "http"

var appName = "bar"

func main() {

	portPtr := flag.Int("port", 8080, "port to listen on")
	namePtr := flag.String("name", "svc", "register on discovery under this name")
	svcName := flag.String("service", "barbaz", "register this service")
	flag.Parse()

	appName = *namePtr

	_ = configureDiscovery(*svcName, *portPtr)
	r := configureRouter()
	go func() {
		_ = http.ListenAndServe(
			fmt.Sprintf("%s:%d", serviceHost, *portPtr),
			r,
		)
	}()
	common.HandleSignals(syscall.SIGINT, syscall.SIGTERM)
}

func configureDiscovery(svcName string, port int) disco.DiscoRegistry {
	for {
		cfg := disco.
			DefaultConfig().
			SkipSslVerify(). // to use with self-signed certificate on disco service
			WithName(svcName).
			WithEndpoints([]string{fmt.Sprintf("%s://%s:%d", serviceProto, serviceHost, port)}).
			WithRetry(2, 2*time.Second)
		//WithTimeout(5 * time.Second).
		clnt, err := disco.NewDiscoHttpClient(cfg)
		if err != nil {
			logger.Warning("join error: %s", err.Error())
			time.Sleep(time.Second)
			continue
		}
		return clnt.Registry()
	}
}
func configureRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/svc", handleBar).Methods("GET")
	return router
}
func handleBar(w http.ResponseWriter, r *http.Request) {
	common.WriteResponseMessage(w, http.StatusOK, http.StatusOK, "message", fmt.Sprintf("[%s] %s", appName, strings.ToUpper(appName)))
}
