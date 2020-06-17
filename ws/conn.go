package ws

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type UnitConn struct {
	wsConnect *websocket.Conn
	inChan    chan []byte // receive buf
	outChan   chan []byte // send buf

	device          string
	addr            string
	lastHeartbeatTs int64

	state     StateType
	closeChan chan byte
	mutex     sync.Mutex // 对closeChan关闭上锁
}

func InitUnitConn(wsConn *websocket.Conn) (c *UnitConn) {
	c = &UnitConn{
		wsConnect: wsConn,
		state:     STATE_INIT,
		inChan:    make(chan []byte, 256), // kafka
		outChan:   make(chan []byte, 256),
		addr:      wsConn.RemoteAddr().String(),
		closeChan: make(chan byte, 1),
	}
	go c.writeLoop()
	go c.readLoop()
	return
}

func (c *UnitConn) String() string {
	return fmt.Sprintf("addr: %s, state: %d", c.addr, c.state)
}

func (c *UnitConn) readLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = c.wsConnect.ReadMessage(); err != nil {
			fmt.Println("ReadMessage err", err)
			c.Close()
		}
		select {
		case c.inChan <- data:
		case <-c.closeChan:
			fmt.Println("UnitConn is closed, readLoop exit!")
		}
	}
}

func (c *UnitConn) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-c.outChan:
		case <-c.closeChan:
			fmt.Println("UnitConn is closed, writeLoop exit!")
			return
		}
		if err = c.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Println("WriteMessage err", err)
			c.Close()
		}
	}
}

// Close _
func (c *UnitConn) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 线程安全，可多次调用
	if err := c.wsConnect.Close(); err != nil {
		fmt.Printf("client %s, Close net.Conn err: %s\n", c, err)
	}

	if c.state != STATE_CLOSED {
		close(c.closeChan)
		c.state = STATE_CLOSED
		fmt.Println("UnitConn is closed")
	}
}

// ReadMessage 消费收到的消息
func (c *UnitConn) ReadMessage() (data []byte, err error) {
	select {
	case data = <-c.inChan:
		fmt.Println("receive:", string(data))
	case <-c.closeChan:
		err = errors.New("UnitConn is closed")
	}
	return
}

// WriteMessage 发消息到队列
func (c *UnitConn) WriteMessage(data []byte) (err error) {
	select {
	case c.outChan <- data:
	case <-c.closeChan:
		err = errors.New("UnitConn is closed")
	}
	return
}

func (c *UnitConn) chgCLOSED2TCP_READY(conn *websocket.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.state == STATE_TCP_READY {
		fmt.Printf("client %s, is already: %s\n", c, STATE_TCP_READY)
		return
	}

	c.addr = conn.RemoteAddr().String()
	// c.msgs = make(map[uint64]*AckMessage)
	c.lastHeartbeatTs = time.Now().Unix()

	old := c.state
	c.state = STATE_TCP_READY
	fmt.Printf("client %s, change %s to %s\n", c, old, c.state)
}

func (c *UnitConn) chgESTABLISHED(pkt interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.state == STATE_ESTABLISHED {
		// c.sendSYNACK(c.cid)
		return false
	}

	old := c.state
	c.state = STATE_SYN_RECEIVED
	fmt.Printf("client %s, state: %s to %s\n", c, old, c.state)

	// c.sendSYNACK(c.cid)
	c.state = STATE_ESTABLISHED

	// DefaultMgr.addClient(c)
	// DefaultMgr.sendtoRTCSRV(c, pkt)

	return true
}

func (c *UnitConn) chgCLOSED() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.state == STATE_CLOSED {
		fmt.Printf("client %s, is already: %s\n", c, STATE_CLOSED)
		return
	}

	if c.state == STATE_ESTABLISHED {
		// DefaultMgr.delClient(c.cid, c.addr)

		// pkt := &chatpkt.Packet{
		// 	Type:     chatpkt.Packet_CLOSE.Enum(),
		// 	ClientId: proto.String(c.cid),
		// }
		// DefaultMgr.sendtoRTCSRV(c, pkt)
		// DefaultMgr.writeChatLogs("clos", c.uuid, pkt)
	}

	c.Close()

	// for k, v := range c.msgs {
	// 	select {
	// 	case v.ack <- false:
	// 	default:
	// 		logger.Trace("client %s, msgid: %d no wait notify", c, k)
	// 	}
	// }

	old := c.state
	c.state = STATE_CLOSED
	fmt.Printf("client %s, change %s to %s\n", c, old, c.state)
}
