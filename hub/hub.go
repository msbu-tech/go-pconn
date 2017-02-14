package hub

import (
    "log"
)

type Hub struct {

}

func (hub *Hub) ConnExists (cuid string) bool {
    log.Println("i'm in connExists")
    return false
}

func (hub *Hub) AddConn (cuid string, pconn interface{}) error{
    log.Println("i'm in addConn")
    return nil
}

func (hub *Hub) RemoveConn (cuid string) error {
    log.Println("i'm in removeConn")
    return nil
}