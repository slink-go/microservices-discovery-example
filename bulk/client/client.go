package client

import (
	"common"
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

type Client struct {
	name     string
	service  string
	port     uint16
	registry disco.DiscoClient
}

func New(app, service string, port uint16) *Client {
	return &Client{
		name:    app,
		service: service,
		port:    port,
	}
}
func (c *Client) Run() {
	go func() {
		_ = c.configureDiscovery(c.service, int(c.port))
		r := c.configureRouter()
		go func() {
			_ = http.ListenAndServe(fmt.Sprintf("%s:%d", serviceHost, c.port), r)
		}()
		common.HandleSignals(syscall.SIGINT, syscall.SIGTERM)
	}()
}
func (c *Client) leave() {
	_ = c.registry.Leave()
}
func (c *Client) configureDiscovery(svcName string, port int) disco.DiscoRegistry {
	for {
		cfg := disco.
			DefaultConfig().
			SkipSslVerify().
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
		c.registry = clnt
		return clnt.Registry()
	}
}
func (c *Client) configureRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/svc", c.handleSvc).Methods("GET")
	return router
}
func (c *Client) handleSvc(w http.ResponseWriter, r *http.Request) {
	common.WriteResponseMessage(w, http.StatusOK, http.StatusOK,
		"message",
		fmt.Sprintf("[%s] %s", c.name, strings.ToUpper(c.name)),
	)
}
