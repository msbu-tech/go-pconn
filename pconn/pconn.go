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
	timestamp       int64           //建连时间戳
	closed          bool            //连接是否关闭
	connected       bool            //是否已经建立连接（客户端是否已经发送connect命令）
	hub             *MyHub          //连接所在的HUB
	c               *websocket.Conn //链接的ws指针
	connectedTimer  *time.Timer     //连接计时器，用于关闭超时未发送connect的连接
	connTimeoutChan chan bool
	pushChan chan msg.Message
}

//新建连接，一般由Hub接收到连接请求时发起
func New(h *MyHub, conn *websocket.Conn) *Pconn {
	c := Pconn{
		rid:             uuid.NewV4().String(),
		hub:             h,
		c:               conn,
		timestamp:       time.Now().Unix(),
		//test
		cuid : string(time.Now().Unix()),
		//test end
		closed:          false,
		connected:       false,
		connTimeoutChan: make(chan bool),
		pushChan:        make(chan msg.Message),
	}
	log.Printf("got a new connection. cuid: %s", c.cuid)
	go c.run()

	//设置客户端未连接超时，用于清理建连接后未发请求的连接
	if true {
		c.connectedTimer = time.AfterFunc(1000*time.Second, c.connectTimeout)
	}

	return &c
}

//请求主循环，监听ws的消息，和关闭链接的channel
func (c *Pconn) run() {
	var message = msg.Message{}
	for {
		select {
		case <-c.connTimeoutChan:
			if !c.connected && !c.closed {
				c.disconnect()
				log.Printf("connenction timeout, goroutine exit, rid: %s", c.rid)
				return
			}
		case message = <- c.pushChan:
			err := c.push(&message)
			if err != nil {
				//todo 失败处理
				log.Println("push error: ", err)
			}
		default:
			_, message, err := c.c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				c.disconnect()
				return
			}
			//TODO 解包，处理msg
			log.Printf("recv: %s", message)
		}
	}
}

func (c *Pconn) Push(messageStr string) {
	message := msg.Message{Body:messageStr}
	c.pushChan <- message
}

func (c *Pconn) push(message *msg.Message) error {
	messageStr := message.Body
	err := c.c.WriteMessage(websocket.TextMessage, []byte(messageStr))
	if err != nil {
		log.Println("write error:", err)
		return err
	}
	return nil
}

func (c *Pconn) connect() error {
	err := c.hub.AddPconn(c.cuid, c)
	if err != nil {
		log.Println("connect error:", err)
		return err
	}
	c.connected = true
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
	log.Printf("close unconnected link, rid: %s", c.rid)
}
