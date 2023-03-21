package mircat

import (
	"encoding/base64"
	"fmt"
)

type ConnManager struct {
	app      *App
	clients  []*TcpClient
	server   *TCPServer
	transfer *TCPTransfer
	cfg      *Config
}

func NewConnManager(app *App, cfg *Config) *ConnManager {
	return &ConnManager{
		app:      app,
		clients:  []*TcpClient{},
		server:   NewTCPServer(app),
		transfer: NewTCPTransfer(app),
		cfg:      cfg,
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

// TransferTcpStart starts the TCP transfer between the source and destination addresses specified in the configuration file of the connection manager.
// It returns true if the transfer was successfully started, otherwise it returns false. It emits an event with the message "transfer-tcp-info" if the server is successfully listening on the source address.
// It emits an event with the message "transfer-tcp-error" if there is an error in starting the transfer, along with the error message.
//
// Parameters:
// - srcAddress: a string that represents the source IP address and port for the TCP transfer
// - dstAddress: a string that represents the destination IP address and port for the TCP transfer
func (c *ConnManager) TransferTcpStart() bool {
	srcAddress := c.cfg.Transfer.SrcAddr + ":" + c.cfg.Transfer.SrcPort
	dstAddress := c.cfg.Transfer.DstAddr + ":" + c.cfg.Transfer.DstPort
	err := c.transfer.Start(srcAddress, dstAddress)
	if err != nil {
		c.app.EventsEmit("transfer-tcp-error", "server", fmt.Sprintf("failed to listen on %s: %v", srcAddress, err))
		fmt.Printf("Failed to start transfer server : %v\n", err)
		return false
	}
	c.app.EventsEmit("transfer-tcp-info", "server", fmt.Sprintf("listening on %s", srcAddress))
	return true
}

// TransferTcpStop stops the transfer TCP server if it is currently running.
// It returns a boolean indicating whether the server was stopped successfully or not.
func (c *ConnManager) TransferTcpStop() bool {
	if c.transfer == nil {
		return false
	}
	c.transfer.Stop()
	c.app.EventsEmit("transfer-tcp-info", "server", "transfer server stopped")
	return true
}

// TransferSendToServer transfers base64 encoded data to the server via TCP connection.
// It first checks if the transfer server is started and emits an event with an error message if not.
// It then decodes the base64 data and sends it to the server using the transfer object.
// If the base64 decoding fails, it emits an event with an error message.
//
// Parameters:
// - client: a string representing the ID of the client sending the data
// - base64Data: a string representing the base64 encoded data to be sent to the server
func (c *ConnManager) TransferSendToServer(client string, base64Data string) {
	if c.transfer == nil || c.transfer.listener == nil {
		c.app.EventsEmit("transfer-tcp-error", client, "transfer server not started")
		return
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.app.EventsEmit("transfer-tcp-error", client, base64Data+" decode failed")
		return
	}
	c.transfer.SendToServer(client, decodedBytes)
}

// TransferSendToClient transfers the decoded data to a specific client via a transfer server.
// It first checks if the transfer server is running and if not, emits a "transfer-tcp-error" event and returns.
// If the transfer server is running, it decodes the base64 data and sends it to the specified client using the transfer server.
// Parameters:
// - client: a string representing the identifier of the client that will receive the data.
// - base64Data: a string representing the data to be transferred, encoded in base64 format.
func (c *ConnManager) TransferSendToClient(client string, base64Data string) {
	if c.transfer == nil || c.transfer.listener == nil {
		c.app.EventsEmit("transfer-tcp-error", client, "transfer server not started")
		return
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.app.EventsEmit("transfer-tcp-error", client, base64Data+" decode failed")
		return
	}
	c.transfer.SendToClient(client, decodedBytes)
}

// TransferBroadcastToServer transfers a base64 encoded string to the server using the connection manager's transfer object.
// If the transfer object or its listener is not available, it emits a 'transfer-tcp-error' event with the error message and returns.
// If decoding of the base64 encoded string fails, it emits a 'transfer-tcp-error' event with the error message and returns.
// Otherwise, it decodes the base64 string into bytes and broadcasts it to the server via the transfer object.
//
// Parameters:
// - base64Data: A base64 encoded string to be transferred to the server.
func (c *ConnManager) TransferBroadcastToServer(base64Data string) {
	if c.transfer == nil || c.transfer.listener == nil {
		c.app.EventsEmit("transfer-tcp-error", "server", "transfer server not started")
		return
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.app.EventsEmit("transfer-tcp-error", "server", base64Data+" decode failed")
		return
	}
	c.transfer.BroadcastToServer(decodedBytes)
}

// TransferBroadcastToClient transfers a base64 encoded string to the connected clients.
// If the server or listener is not started, a "transfer-tcp-error" event will be emitted with the corresponding error message.
// If the base64 decoding fails, a "transfer-tcp-error" event will be emitted with the corresponding error message.
// Parameters:
// - base64Data: the base64 encoded string to be transferred to the clients.
func (c *ConnManager) TransferBroadcastToClient(base64Data string) {
	if c.server == nil || c.server.listener == nil {
		c.app.EventsEmit("transfer-tcp-error", "server", "transfer server not started")
		return
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.app.EventsEmit("transfer-tcp-error", "server", base64Data+" decode failed")
		return
	}
	c.transfer.BroadcastToClient(decodedBytes)
}
