package pconn

import (
    "log"
    "sync"
    "time"

    "github.com/msbu-tech/go-pconn/hub"
    "github.com/gorilla/websocket"
    "github.com/satori/go.uuid"
)

//连接定义
type Pconn struct {
    sync.RWMutex    //锁
    rid string      //请求id
    cuid string     //uid
    timestamp int64 //建连时间戳
    closed bool     //连接是否关闭
    connected bool  //是否已经建立连接（客户端是否已经发送connect命令）
    hub *hub.Hub    //连接所在的HUB
    c *websocket.Conn   //链接的ws指针
    connectedTimer *time.Timer  //连接计时器，用于关闭超时未发送connect的连接
    closeChan chan bool //关闭连接的channel
}

//新建连接，一般由Hub接收到连接请求时发起
func New(h *hub.Hub, conn *websocket.Conn) *Pconn  {
    c := Pconn{
        rid: uuid.NewV4().String(),
        hub: h,
        c: conn,
        timestamp: time.Now().Unix(),
        closed: false,
        connected: false,
        closeChan: make(chan bool, 1),
    }

    go c.run()

    //设置客户端未连接超时，用于清理建连接后未发请求的连接
    if true{
        c.connectedTimer = time.AfterFunc(10 * time.Second, c.closeUnconnected)
    }

    return &c
}

//请求主循环，监听ws的消息，和关闭链接的channel
func (c *Pconn) run() {
    for {
        select {
        case <- c.closeChan:
            log.Printf("connenction close, goroutine exit, rid: %s", c.rid)
            return
        default:
            _, message, err := c.c.ReadMessage()
            if err != nil {
                log.Println("read:", err)
                return
            }
            //TODO 解包，处理msg
            log.Printf("recv: %s", message)
        }
    }
}

func (c *Pconn) push(message []byte) {
    err := c.c.WriteMessage(websocket.TextMessage, message)
    if err != nil {
        log.Println("write:", err)
        //error
    }
}

func (c *Pconn) connect() error  {
    if !c.hub.ConnExists(c.cuid) {
        err := c.hub.AddConn(c.cuid, c)
        if err != nil{
            log.Println("connect error:", err)
            return err
        }
    }
    return nil
}

func (c *Pconn) disconnect() error {
    err :=c.c.Close()
    if err != nil{
        log.Println("disconnect error:", err)
        return err
    }
    err = c.hub.RemoveConn(c.cuid)
    if err != nil{
        log.Println("disconnect error:", err)
        return err
    }
    log.Printf("disconnect link, rid: %s", c.rid)
    c.closeChan <- true
    return nil
}

func (c *Pconn) closeUnconnected() {
    c.RLock()
    connected := c.connected
    closed := c.closed
    c.RUnlock()
    if !connected && !closed{
        c.disconnect()
        log.Printf("close unconnected link, rid: %s", c.rid)
    }
}
