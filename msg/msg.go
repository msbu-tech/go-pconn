package msg

import (
	"encoding/json"
	"log"
)

type PushMsgReq struct {
	Push_type  int
	Channel_id string
	Device_ids []string
	Content    string
}

type PushMsgRes struct {
	Errno  int
	Errmsg string
}

type PushMsg struct {
	Body string
}

type PconnQueryMsgRes struct {
	Errno       int
	Errmsg      string
	Pconn_Count int
	Device_ids  []string
}

func (p *PushMsg) ToString() string {
	return p.Body
}

func (p *PushMsg) SetContent(content string) {
	p.Body = content
}

type ClientMsg struct {
	Cuid string `json:"cuid"`
	Cmd  string `json:"cmd"`
	Body string `json:"body"`
}

func (c *ClientMsg) GetCuid() string {
	return c.Cuid
}

func (c *ClientMsg) GetCmd() string {
	return c.Cmd
}

func (c *ClientMsg) GetBody() string {
	return c.Body
}

func NewClientMsg(msg string, cMsg *ClientMsg) error {
	err := json.Unmarshal([]byte(msg), cMsg)
	log.Printf("%v", cMsg)
	if err != nil {
		log.Printf("client message parse error, msg:%s, error:%s", msg, err)
		return err
	}
	return nil
}
