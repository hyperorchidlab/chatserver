package httpservice

import (
	"context"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/chatserver/httpservice/api"
	"log"
	"net/http"
	"strconv"
	"time"
)

var webserver *http.Server


func StartWebDaemon() {
	mux := http.NewServeMux()

	mux.Handle("/ajax/userreg", &api.UserRegister{})
	mux.Handle("/ajax/cmd",&api.MessageDispatch{})

	addr := ":" + strconv.Itoa(config.GetCSC().MgtHttpPort)

	log.Println("Web Server Start at", addr)

	webserver = &http.Server{Addr: addr, Handler: mux}

	log.Fatal(webserver.ListenAndServe())

}

func StopWebDaemon() {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	webserver.Shutdown(ctx)

}
