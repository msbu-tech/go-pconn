package server

import (
	//"fmt"
	"log"
	//"time"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"

	//"github.com/msbu-tech/go-pconn/hub"
	"github.com/msbu-tech/go-pconn/msg"
	"github.com/msbu-tech/go-pconn/pconn"
)

var (
	upgrader  websocket.Upgrader
	h         *pconn.MyHub
	conn_addr string
	push_addr string
)

func init() {
	upgrader = websocket.Upgrader{}
	h = pconn.NewHub()
	conn_addr = ":8077"
	push_addr = ":8078"

}

func StartPconnSrv() error {
	log.Println("start connSrv...")
	log.Println("connSrv listen at 127.0.0.1:8077...")
	go h.Run()
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(conn_addr, nil)
	if err != nil {
		log.Fatal("StartPconnSrv ListenAndServe: ", err)
		return errors.New("StartPconnSrv ListenAndServe failed...")
	}
	return nil
}

func StartPusherSrv() error {
	log.Println("start pusherSrv...")
	log.Println("pusherSrv listen at 127.0.0.1:8078...")
	http.HandleFunc("/push", servePush)
	err := http.ListenAndServe(push_addr, nil)
	if err != nil {
		log.Fatal("StartPusherSrv ListenAndServe failed...")
		return errors.New("StartPusherSrv failed...")
	}
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

func servePush(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var push_req msg.PushMsgReq
	err = json.Unmarshal(body, &push_req)
	if err != nil {
		panic(err)
	}
	log.Printf(push_req.Content)
	for _, device_id := range push_req.Device_ids {
		c := h.GetPconn(device_id)
		if c == nil {
			log.Println("Conn is not exists, device_id: ", device_id)
			continue
		}
		c.Push(push_req.Content)
	}
	push_res := msg.PushMsgRes{
		Errno:  0,
		Errmsg: "Success",
	}
	b, err := json.Marshal(push_res)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/json")
	w.Write(b)
}
