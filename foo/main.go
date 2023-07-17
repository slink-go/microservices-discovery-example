package main

import (
	"common"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	disco "github.com/slink-go/disco-go"
	"github.com/slink-go/httpclient"
	"net/http"
	"syscall"
	"time"
)

const servicePort = 8083
const serviceHost = "localhost"
const serviceProto = "http"
const serviceName = "foo"
const remoteServiceName = "barbaz"

var client disco.LoadBalancingHttpClient
var registry disco.DiscoRegistry

func main() {
	client = configureHttpClient()
	r := configureRouter()
	go func() {
		_ = http.ListenAndServe(
			fmt.Sprintf("%s:%d", serviceHost, servicePort),
			r,
		)
	}()
	common.HandleSignals(syscall.SIGINT, syscall.SIGTERM)
}

func configureDiscovery() disco.DiscoRegistry {
	cfg := disco.DefaultConfig().
		SkipSslVerify(). // to use with self-signed certificate on disco service
		//WithToken(os.Getenv("DISCO_TOKEN")).
		WithName(serviceName).
		//WithDisco([]string{os.Getenv("DISCO_URL")}).
		//WithEndpoints([]string{fmt.Sprintf("%s://%s:%d", serviceProto, serviceHost, servicePort)})
		//WithBreaker(2).
		WithRetry(2, 2*time.Second)
	//WithTimeout(5 * time.Second)
	registry = disco.NewDiscoHttpClientPanicOnAuthError(cfg).Registry()
	return registry
}
func configureRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/foo", handleFoo).Methods("GET")
	return router
}
func configureHttpClient() disco.LoadBalancingHttpClient {
	return disco.NewLbClient(
		configureDiscovery(),
		httpclient.New().WithNoAuth().
			//WithBreaker(3).
			WithRetry(3, time.Second),
	)
}

func handleFoo(w http.ResponseWriter, r *http.Request) {
	//b, code, err := client.Get("http://localhost:8081/api/bar")
	b, _, code, err := client.Get(fmt.Sprintf("%s://%s/api/svc", serviceProto, remoteServiceName), nil)
	if err != nil {
		common.WriteResponseError(w, code, err)
		return
	}
	if code >= 400 {
		common.WriteResponseError(w, code, err)
		return
	}
	b, code = processResponse(b)
	common.WriteResponseBytes(w, code, b)
}
func processResponse(input []byte) ([]byte, int) {

	var v map[string]any
	err := json.Unmarshal(input, &v)
	if err != nil {
		return []byte(err.Error()), http.StatusInternalServerError
	}

	message, ok := v["message"]
	if ok {
		s := fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), message)
		m := map[string]any{}
		m["message"] = s
		b, er := json.Marshal(m)
		if er != nil {
			return []byte(er.Error()), http.StatusInternalServerError
		}
		return b, http.StatusOK
	}

	errmsg, ok := v["error"]
	if ok {
		s := fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), errmsg)
		m := map[string]any{}
		m["error"] = s
		b, er := json.Marshal(m)
		if er != nil {
			return []byte(er.Error()), http.StatusInternalServerError
		}
		return b, http.StatusInternalServerError
	}

	return input, http.StatusOK
}
