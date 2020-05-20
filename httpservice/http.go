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

var(
	register string = "userreg"
	addfriend string = "addfriend"
	removefriend string = "removefriend"
	addgroup string = "addgroup"
	removegroup string = "removegroup"
	joingroup string = "joingroup"
	quitgroup string = "quitgroup"

	groupmsg string = "groupmsg"
	p2pmsg string = "p2pmsg"
)

func StartWebDaemon() {
	mux := http.NewServeMux()

	mux.Handle("/ajax/userreg", &api.UserRegister{})

	addr := ":" + strconv.Itoa(config.GetCSC().MgtHttpPort)

	log.Println("Web Server Start at", addr)

	webserver = &http.Server{Addr: addr, Handler: mux}

	log.Fatal(webserver.ListenAndServe())

}

func StopWebDaemon() {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	webserver.Shutdown(ctx)

}
