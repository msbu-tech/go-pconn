package server

import (
	//"fmt"
	"log"
	//"time"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/msbu-tech/go-pconn/hub"
	"github.com/msbu-tech/go-pconn/pconn"
)

var (
	upgrader websocket.Upgrader
	h        *hub.Hub
	addr     string
)

func init() {
	upgrader = websocket.Upgrader{}
	h = hub.NewHub()
	addr = ":8079"
}

func StartPconnSrv() {
	go h.Run()
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("start connSrv success...")
	log.Println("connSrv listen at 127.0.0.1:8089...")
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade:", err)
		return
	}
	pconn.New(h, c)
}
