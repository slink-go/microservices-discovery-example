package common

import (
	"fmt"
	"github.com/slink-go/httpclient"
	"github.com/slink-go/logger"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func WriteResponseStr(w http.ResponseWriter, code int, str string) {
	WriteResponseBytes(w, code, []byte(fmt.Sprintf("%s\n", str)))
}
func WriteResponseBytes(w http.ResponseWriter, code int, data []byte) {
	//logger.Notice("response: %s", strings.TrimSpace(string(data)))
	w.WriteHeader(code)
	_, err := w.Write(data)
	if err != nil {
		return
	}
}
func WriteResponseMessage(w http.ResponseWriter, code, ecode int, key, value string) {
	WriteResponseStr(w, code, fmt.Sprintf("{\"%s\": \"%s\",\"code\": %d}", key, value, ecode))
}
func WriteResponseError(w http.ResponseWriter, code int, err error) {
	ecode := code
	switch err.(type) {
	case *httpclient.HttpError:
		ecode = err.(*httpclient.HttpError).Code()
	case *httpclient.ServiceUnreachableError:
		ecode = err.(*httpclient.ServiceUnreachableError).Code()
	case *httpclient.ConnectionRefusedError:
		ecode = err.(*httpclient.ConnectionRefusedError).Code()
	}
	WriteResponseMessage(w, code, ecode, "error", strings.ReplaceAll(err.Error(), "\"", "\\\""))
}
func HandleSignals(signals ...os.Signal) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, signals...)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		logger.Debug("[signal] received %s signal", sig)
		done <- true
	}()
	<-done
	time.Sleep(500 * time.Millisecond) // let disco client to send "leave" message
}
