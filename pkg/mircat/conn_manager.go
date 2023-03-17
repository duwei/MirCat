package mircat

import (
	"encoding/base64"
	"fmt"
)

type ConnManager struct {
	app     *App
	clients []*TcpClient
	server  *TCPServer
	cfg     *Config
}

func NewConnManager(app *App, cfg *Config) *ConnManager {
	return &ConnManager{
		app:     app,
		clients: []*TcpClient{},
		server:  NewTCPServer(app),
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

// ServerTcpStart starts the TCP server for the connection manager.
// It takes the server's address from the configuration file, starts the server, and returns a boolean indicating success or failure.
// If the server fails to start, it emits a "server-tcp-error" event with the error message and returns false.
// If the server starts successfully, it emits a "server-tcp-info" event with the server's address and returns true.
func (c *ConnManager) ServerTcpStart() bool {
	address := c.cfg.Server.TcpAddr + ":" + c.cfg.Server.TcpPort
	err := c.server.Start(address)
	if err != nil {
		c.app.EventsEmit("server-tcp-error", "server", fmt.Sprintf("failed to listen on %s: %v", address, err))
		fmt.Printf("Failed to start tcp server : %v\n", err)
		return false
	}
	c.app.EventsEmit("server-tcp-info", "server", fmt.Sprintf("listening on %s", address))
	return true
}

// ServerTcpStop stops the TCP server for the connection manager.
// It checks if the server is running and stops it if it is, then emits a "server-tcp-info" event indicating the server has stopped and returns true.
// If the server is not running, it returns false.
func (c *ConnManager) ServerTcpStop() bool {
	if c.server == nil {
		return false
	}
	c.server.Stop()
	c.app.EventsEmit("server-tcp-info", "server", "tcp server stopped")
	return true
}

// ServerSendMessage sends a message to a specific client over TCP connection.
//
// Parameters:
// - client: the identifier of the target client.
// - base64Data: the message data encoded in base64 format.
//
// If the server is not started or the listener is not initialized, an error event will be emitted
// through the app instance.
// If the base64Data is not a valid base64-encoded string, an error event will also be emitted.
func (c *ConnManager) ServerSendMessage(client string, base64Data string) {
	if c.server == nil || c.server.listener == nil {
		c.app.EventsEmit("server-tcp-error", client, "server not started")
		return
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.app.EventsEmit("server-tcp-error", client, base64Data+" decode failed")
		return
	}
	c.server.SendMessage(client, decodedBytes)
}

// ServerBroadcastMessage broadcasts a message to all connected clients over TCP connection.
//
// Parameters:
// - base64Data: the message data encoded in base64 format.
//
// If the server is not started or the listener is not initialized, an error event will be emitted
// through the app instance.
// If the base64Data is not a valid base64-encoded string, an error event will also be emitted.
func (c *ConnManager) ServerBroadcastMessage(base64Data string) {
	if c.server == nil || c.server.listener == nil {
		c.app.EventsEmit("server-tcp-error", "server", "server not started")
		return
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.app.EventsEmit("server-tcp-error", "server", base64Data+" decode failed")
		return
	}
	c.server.BroadcastMessage(decodedBytes)
}
