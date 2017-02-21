package pconn

import (
//"fmt"
//"log"
//"errors"
//"github.com/msbu-tech/go-pconn/pconn"
)

type AddPconnMsg struct {
	cuid  string
	pconn *Pconn
}

type DelPconnMsg struct {
	cuid string
}

type MyHub struct {
	pconn_pool map[string]*Pconn
	add_conn   chan *AddPconnMsg
	del_conn   chan *DelPconnMsg
}

func (h *MyHub) AddPconn(cuid string, pconn *Pconn) error {
	addPconnMsg := &AddPconnMsg{
		cuid:  cuid,
		pconn: pconn,
	}
	h.add_conn <- addPconnMsg
	return nil
}

func (h *MyHub) DelPconn(cuid string) error {
	delPconnMsg := &DelPconnMsg{
		cuid: cuid,
	}
	h.del_conn <- delPconnMsg
	return nil
}

func (h *MyHub) IsPconnExist(cuid string) bool {
	if _, ok := h.pconn_pool[cuid]; ok {
		return true
	}
	return false
}

func (h *MyHub) GetPconn(cuid string) *Pconn {
	if _, ok := h.pconn_pool[cuid]; ok {
		return h.pconn_pool[cuid]
	}
	return nil
}

func NewHub() *MyHub {
	return &MyHub{
		pconn_pool: make(map[string]*Pconn),
		add_conn:   make(chan *AddPconnMsg),
		del_conn:   make(chan *DelPconnMsg),
	}
}

func (h *MyHub) Run() {
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
