package mircat

import (
	"encoding/base64"
	"fmt"
)

type ConnManager struct {
	app     *App
	clients []*TcpClient
	cfg     *Config
}

func NewConnManager(app *App, cfg *Config) *ConnManager {
	return &ConnManager{
		app:     app,
		clients: []*TcpClient{},
		cfg:     cfg,
	}
}

// ClientTcpOpen opens a new TCP client connection and returns its index.
// Returns:
// - int: the index of the newly opened client connection.
func (c *ConnManager) ClientTcpOpen() int {
	tcpClient, err := NewTcpClient(c.cfg.Client.ServerIp+":"+c.cfg.Client.ServerPort, c.app)
	if err != nil {
		c.app.EventsEmit("client-tcp-error", -1, fmt.Sprintf("%v", err))
		fmt.Printf("Failed to connect: %v\n", err)
		return -1
	}
	if c.clients == nil {
		c.clients = []*TcpClient{}
	}
	c.clients = append(c.clients, tcpClient)
	tcpClient.index = len(c.clients) - 1
	c.app.EventsEmit("client-tcp-info", tcpClient.index, "connection opened")
	return tcpClient.index
}

// ClientTcpSend sends data to a specified TCP client.
// Parameters:
// - index (int): the index of the TCP client.
// - base64Data (string): the data to send, encoded in base64 format.
func (c *ConnManager) ClientTcpSend(index int, base64Data string) {
	if c.clients == nil || index >= len(c.clients) || index < 0 {
		c.app.EventsEmit("client-tcp-error", index, "invalid client index")
		return
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.app.EventsEmit("client-tcp-error", index, base64Data+" decode failed")
		return
	}
	c.clients[index].Send(decodedBytes)
}

// ClientTcpClose closes a specified TCP client connection.
// Parameters:
// - index (int): the index of the TCP client to close.
func (c *ConnManager) ClientTcpClose(index int) {
	if c.clients == nil || index >= len(c.clients) || index < 0 {
		c.app.EventsEmit("client-tcp-error", index, "invalid client index")
		return
	}
	c.clients[index].Shutdown()
}

// ClientTcpCloseAll closes all currently open TCP client connections.
func (c *ConnManager) ClientTcpCloseAll() {
	for _, client := range c.clients {
		client.Shutdown()
	}
	c.clients = []*TcpClient{}
}
