package server

import (
	//"fmt"
	"log"
	//"time"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"

	//"github.com/msbu-tech/go-pconn/hub"
	"github.com/msbu-tech/go-pconn/pconn"
)

var (
	upgrader websocket.Upgrader
	h        *pconn.MyHub
	addr     string
)

func init() {
	upgrader = websocket.Upgrader{}
	h = pconn.NewHub()
	addr = ":8079"
}

func StartPconnSrv() error {
	go h.Run()
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return errors.New("ListenAndServe failed...")
	}
	log.Println("start connSrv success...")
	log.Println("connSrv listen at 127.0.0.1:8089...")
	return nil
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade:", err)
		return
	}
	pconn.New(h, c)
}
