package mircat

import (
	"fmt"
	"net"
	"time"
)

const RECONNECT_INTERVAL = time.Second

type TcpClient struct {
	address    string      // 连接的地址
	conn       net.Conn    // 实际的网络连接对象
	sendChan   chan []byte // 发送数据的通道
	recvChan   chan []byte // 接收数据的通道
	isShutdown bool        // 是否关闭
	index      int
	app        *App
}

func NewTcpClient(address string, app *App) (*TcpClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	c := &TcpClient{
		address:    address,
		conn:       conn,
		sendChan:   make(chan []byte),
		recvChan:   make(chan []byte),
		isShutdown: false,
		index:      -1,
		app:        app,
	}
	go c.startSending()
	go c.startReceiving()
	return c, nil
}

func (c *TcpClient) startSending() {
	for {
		select {
		case data := <-c.sendChan:
			_, err := c.conn.Write(data)
			if err != nil {
				if c.isShutdown {
					return
				}
				c.reconnect()
				return
			}
		case <-time.After(time.Second):
			// do nothing
		}
	}
}

func (c *TcpClient) startReceiving() {
	buffer := make([]byte, 4096)
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			if c.isShutdown {
				return
			}
			c.app.EventsEmit("client-tcp-error", c.index, fmt.Sprintf("connection closed: %v", err))
			c.reconnect()
			return
		}
		dst := make([]byte, n)
		copy(dst, buffer[:n])
		c.app.EventsEmit("client-tcp-data", c.index, dst)
		fmt.Printf("Recv data: %v\n", buffer[:n])
		//c.recvChan <- buffer[:n]
	}
}

func (c *TcpClient) Send(data []byte) {
	if c.isShutdown {
		c.app.EventsEmit("client-tcp-error", c.index, "connection closed")
		return
	}
	c.sendChan <- data
}

//func (c *TcpClient) Recv() ([]byte, error) {
//	if c.isShutdown {
//		return nil, fmt.Errorf("TcpClient is shutdown")
//	}
//	select {
//	case data := <-c.recvChan:
//		return data, nil
//	case err := <-c.errorChan:
//		return nil, err
//	}
//}

func (c *TcpClient) Shutdown() {
	if !c.isShutdown {
		c.isShutdown = true
		c.conn.Close()
		close(c.sendChan)
		close(c.recvChan)
		c.app.EventsEmit("client-tcp-info", c.index, "connection closed")
	}
}

func (c *TcpClient) reconnect() {
	for {
		if c.isShutdown {
			return
		}
		conn, err := net.Dial("tcp", c.address)
		if err == nil {
			c.conn = conn
			go c.startSending()
			go c.startReceiving()
			c.app.EventsEmit("client-tcp-info", c.index, "connection reconnected")
			return
		}
		c.app.EventsEmit("client-tcp-info", c.index, "trying to reconnect...")
		time.Sleep(RECONNECT_INTERVAL)
	}
}

//func main() {
//	conn, err := NewTcpClient("127.0.0.1:9000")
//	if err != nil {
//		fmt.Printf("Failed to connect: %v\n", err)
//		return
//	}
//	defer conn.Shutdown()
//
//	// 发送数据
//	err = conn.Send([]byte("hello"))
//	if err != nil {
//		fmt.Printf("Failed to send data: %v\n", err)
//		return
//	}
//
//	for i := 0; i < 10; i++ {
//		// 接收数据
//		data, err := conn.Recv()
//		if err != nil {
//			fmt.Printf("Failed to recv data: %v\n", err)
//			continue
//		}
//		fmt.Printf("Recv data: %v\n", data)
//	}
//
//}
