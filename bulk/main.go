package main

import (
	"bulk/client"
	"flag"
	"github.com/slink-go/logger"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const serviceName = "barbaz"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
var clients = make(map[string]*client.Client)

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	amount := flag.Int("n", 1, "number of clients to start")
	svcName := flag.String("s", serviceName, "service name")
	flag.Parse()
	port := 20001
	for {
		time.Sleep(time.Duration(rand.Intn(3000) + 50))
		if len(clients) == *amount {
			break
		}
		key := randStringRunes(3)
		_, ok := clients[key]
		if ok {
			continue
		}
		clients[key] = client.New(key, *svcName, uint16(port))
		clients[key].Run()
		port++
	}
	handleSignals(syscall.SIGINT, syscall.SIGTERM)
}

func handleSignals(signals ...os.Signal) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, signals...)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		logger.Debug("[signal] received %s signal", sig)
		done <- true
	}()
	<-done
	//for _, c := range clients {
	//	go c.leave()
	//}
	time.Sleep(5000 * time.Millisecond) // let disco client to send "leave" message
}
