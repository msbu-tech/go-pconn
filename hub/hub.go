package hub

import (
	//"fmt"
	//"log"

	"github.com/msbu-tech/go-pconn/pconn"
)

type AddPconnMsg struct {
	cuid  string
	pconn *pconn.Pconn
}

type DelPconnMsg struct {
	cuid string
}

type Hub struct {
	pconn_pool map[string]*pconn.Pconn
	add_conn   chan *AddPconnMsg
	del_conn   chan *DelPconnMsg
}

func (h *Hub) AddPconn(cuid string, pconn *pconn.Pconn) {
	addPconnMsg := &AddPconnMsg{
		cuid:  cuid,
		pconn: pconn,
	}
	h.add_conn <- addPconnMsg
}

func (h *Hub) DelPconn(cuid string) {
	delPconnMsg := &DelPconnMsg{
		cuid: cuid,
	}
	h.del_conn <- delPconnMsg
}

func (h *Hub) IsPconnExist(cuid string) bool {
	if _, ok := h.pconn_pool[cuid]; ok {
		return true
	}
	return false
}

func (h *Hub) GetPconn(cuid string) *pconn.Pconn {
	if _, ok := h.pconn_pool[cuid]; ok {
		return h.pconn_pool[cuid]
	}
	return nil
}

func NewHub() *Hub {
	return &Hub{
		pconn_pool: make(map[string]*pconn.Pconn),
		add_conn:   make(chan *AddPconnMsg),
		del_conn:   make(chan *DelPconnMsg),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case addPconnMsg := <-h.add_conn:
			h.pconn_pool[addPconnMsg.cuid] = addPconnMsg.pconn
		case delPconnMsg := <-h.del_conn:
			if _, ok := h.pconn_pool[delPconnMsg.cuid]; ok {
				delete(h.pconn_pool, delPconnMsg.cuid)
			}
		}
	}
}
