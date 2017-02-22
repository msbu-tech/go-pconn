package main

import (
	"flag"
	"log"
	//"net/http"

	//"github.com/msbu-tech/go-pconn/hub"
	//"github.com/msbu-tech/go-pconn/pconn"
	"github.com/msbu-tech/go-pconn/server"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

//var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options

func main() {
	flag.Parse()
	//log.Printf("addr: %s", *addr)
	log.SetFlags(0)
	//http.HandleFunc("/", handler)
	//log.Fatal(server.StartPconnSrv())
	//log.Fatal(server.StartPusherSrv())
	go server.StartPconnSrv()
	server.StartPusherSrv()
}
