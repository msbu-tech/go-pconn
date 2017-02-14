package main

import (
    "log"
    "flag"
    "net/http"

    "github.com/msbu-tech/go-pconn/hub"
    "github.com/msbu-tech/go-pconn/pconn"
    "github.com/gorilla/websocket"
)

var h = new(hub.Hub)

func handler(w http.ResponseWriter, r *http.Request)  {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    log.Println("connect upgrade success")
    pconn.New(h, c)
}

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true },} // use default options

func main()  {
    flag.Parse()
    log.Printf("addr: %s", *addr)
    log.SetFlags(0)
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(*addr, nil))
}