package pconn

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/msbu-tech/go-pconn/msg"
	"github.com/satori/go.uuid"
)

//连接定义
type Pconn struct {
	rid             string          //请求id
	cuid            string          //uid
	uuid            string          //server端生成的唯一id
	timestamp       int64           //建连时间戳
	closed          bool            //连接是否关闭
	connected       bool            //是否已经建立连接（客户端是否已经发送connect命令）
	hub             *MyHub          //连接所在的HUB
	c               *websocket.Conn //链接的ws指针
	connectedTimer  *time.Timer     //连接计时器，用于关闭超时未发送connect的连接
	connTimeoutChan chan bool
	pushChan chan *msg.PushMsg
	closeChan chan bool
	cmdChan chan *msg.ClientMsg
}

//新建连接，一般由Hub接收到连接请求时发起
func New(h *MyHub, conn *websocket.Conn) *Pconn {
	c := Pconn{
		rid:             uuid.NewV4().String(),
		hub:             h,
		c:               conn,
		timestamp:       time.Now().Unix(),
		uuid : uuid.NewV1().String(),
		closed:          false,
		connected:       false,
		connTimeoutChan: make(chan bool),
		pushChan:        make(chan *msg.PushMsg),
		closeChan:       make(chan bool),
		cmdChan:         make(chan *msg.ClientMsg),
	}
	go c.read()
	go c.processChan()

	//设置客户端未连接超时，用于清理建连接后未发请求的连接
	if true {
		c.connectedTimer = time.AfterFunc(10*time.Second, c.connectTimeout)
	}

	return &c
}

//请求主循环，监听ws的消息
func (c *Pconn) read() {
	for {
		_, message, err := c.c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			c.closeChan <- true
            break
		}
		log.Printf("recv: %s", message)
		clientMsg := new(msg.ClientMsg)
		err = msg.NewClientMsg(string(message), clientMsg)
		if err != nil{
			log.Printf("client message parse error, msg: %s", clientMsg)
			continue
		}
		err = c.dispatchCmd(clientMsg)
		if err != nil{
			log.Printf("client message dispatch error, error: %s", err)
		}
	}
}

func (c *Pconn) dispatchCmd(cMsg *msg.ClientMsg) error  {
	switch cMsg.GetCmd() {
	case "connect":
		c.connectCmd(cMsg)
	case "disconnect":
		c.disconnectCmd(cMsg)
	case "message":
		c.messageCmd(cMsg)
	}
	return nil
}

func (c *Pconn) connectCmd(cMsg *msg.ClientMsg)  {
	c.cmdChan <- cMsg
}

func (c *Pconn) disconnectCmd(cMsg *msg.ClientMsg)  {
	c.cmdChan <- cMsg
}

func (c *Pconn) messageCmd(cMsg *msg.ClientMsg)  {
	c.cmdChan <- cMsg
}

func (c *Pconn) processChan()  {
	for {
		select {
		case <-c.connTimeoutChan:
			if !c.connected && !c.closed {
				c.disconnect()
				log.Printf("connenction timeout, goroutine exit, rid: %s", c.rid)
				return
			}
		case <- c.closeChan:
			c.disconnect()
			return
		case message := <- c.pushChan:
			c.push(message)
		case cMsg := <- c.cmdChan:
			switch cMsg.GetCmd() {
			case Connect:
				c.connect(cMsg.GetCuid())
			case Disconnect:
				c.disconnect()
                return
			case Message:
                //TODO 解析msg，连接后端，转发请求
				log.Printf("got client message: %s", cMsg.GetBody())
            default:
                log.Printf("unknown cmd: %s", cMsg.GetCmd())
			}
		}
	}
}

func (c *Pconn) Push(messageStr string) {
	message := new(msg.PushMsg)
	message.SetContent(messageStr)
	c.pushChan <- message
}

func (c *Pconn) push(message *msg.PushMsg) error {
	messageStr := message.ToString()
	err := c.c.WriteMessage(websocket.TextMessage, []byte(messageStr))
	if err != nil {
		log.Println("write error:", err)
		return err
	}
	return nil
}

func (c *Pconn) connect(cuid string) error {
	if !c.connected{
		c.cuid = cuid
		err := c.hub.AddPconn(c.cuid, c)
		if err != nil {
			log.Println("connect error:", err)
			return err
		}
		log.Printf("conn %s connected", c.cuid)
		c.connected = true
	}
	return nil
}

func (c *Pconn) disconnect() error {
	err := c.c.Close()
	log.Println("close ws connnection")
	if err != nil {
		log.Println("disconnect error[ws close]:", err)
		return err
	}
	err = c.hub.DelPconn(c.cuid)
	if err != nil {
		log.Println("disconnect error[RemoveConn]:", err)
		return err
	}
	c.closed = true

	log.Printf("disconnect link, rid: %s", c.rid)
	return nil
}

func (c *Pconn) connectTimeout() {
	c.connTimeoutChan <- true
}
