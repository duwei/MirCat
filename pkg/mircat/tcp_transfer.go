package mircat

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type TransferConn struct {
	clientConn net.Conn
	serverConn net.Conn
}

type TCPTransfer struct {
	srcAddress      string
	dstAddress      string
	listener        net.Listener
	clients         map[string]TransferConn
	mutex           sync.RWMutex
	broadcastServer chan []byte
	broadcastClient chan []byte
	addClient       chan net.Conn
	removeClient    chan net.Conn
	shutdown        chan bool
	app             *App
}

func NewTCPTransfer(app *App) *TCPTransfer {
	return &TCPTransfer{
		clients:         make(map[string]TransferConn),
		broadcastServer: make(chan []byte),
		broadcastClient: make(chan []byte),
		addClient:       make(chan net.Conn),
		removeClient:    make(chan net.Conn),
		shutdown:        make(chan bool),
		app:             app,
	}
}

func (s *TCPTransfer) Start(srcAddress string, dstAddress string) error {
	var err error
	s.srcAddress = srcAddress
	s.dstAddress = dstAddress
	s.listener, err = net.Listen("tcp", s.srcAddress)
	if err != nil {
		return err
	}
	fmt.Printf("Listening on %s\n", s.srcAddress)

	go s.handleEvents()

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					s.app.EventsEmit("transfer-tcp-error", "server", fmt.Sprintf("error accepting connection: %v", err))
				}
				fmt.Printf("Error accepting connection: %s\n", err.Error())
				return
			}
			s.app.EventsEmit("transfer-tcp-info", conn.RemoteAddr().String(), fmt.Sprintf("client connected: %s", conn.RemoteAddr()))
			fmt.Printf("New client connected: %s\n", conn.RemoteAddr())

			s.addClient <- conn
		}
	}()
	return nil
}

func (s *TCPTransfer) Stop() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
		s.shutdown <- true
	}
}

func (s *TCPTransfer) handleClientConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.app.EventsEmit("transfer-tcp-info", conn.RemoteAddr().String(), fmt.Sprintf("client disconnected: %s", conn.RemoteAddr()))
		fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())

		s.removeClient <- conn
	}()

	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF && !errors.Is(err, net.ErrClosed) {
				s.app.EventsEmit("transfer-tcp-error", conn.RemoteAddr().String(), fmt.Sprintf("error reading from client %s : %v", conn.RemoteAddr(), err))
			}
			fmt.Printf("Error reading from client %s: %s\n", conn.RemoteAddr(), err.Error())
			return
		}
		message := append([]byte{}, buffer[:n]...)
		s.app.EventsEmit("transfer-src-data", conn.RemoteAddr().String(), message)
		//s.broadcast <- message
	}
}

func (s *TCPTransfer) handleServerConnection(clientConn net.Conn, serverConn net.Conn) {
	clientKey := clientConn.RemoteAddr().String()

	defer func() {
		serverConn.Close()
		s.app.EventsEmit("transfer-tcp-info", clientKey, fmt.Sprintf("dst disconnected: %s", serverConn.RemoteAddr()))
		fmt.Printf("Dst client disconnected: %s\n", serverConn.RemoteAddr())
	}()

	buffer := make([]byte, 4096)
	for {
		n, err := serverConn.Read(buffer)
		if err != nil {
			if err != io.EOF && !errors.Is(err, net.ErrClosed) {
				s.app.EventsEmit("transfer-tcp-error", clientKey, fmt.Sprintf("error reading from client %s : %v", serverConn.RemoteAddr(), err))
			}
			fmt.Printf("Error reading from client %s: %s\n", serverConn.RemoteAddr(), err.Error())
			if s.getTransferConn(clientKey) == nil {
				return
			}
			s.reconnect(clientKey)
			continue
		}
		message := append([]byte{}, buffer[:n]...)
		s.app.EventsEmit("transfer-dst-data", clientConn.RemoteAddr().String(), message)
		//s.broadcast <- message
	}
}

func (s *TCPTransfer) getTransferConn(clientKey string) *TransferConn {
	s.mutex.RLock()
	transferConn, ok := s.clients[clientKey]
	s.mutex.RUnlock()
	if !ok {
		return nil
	}
	return &transferConn
}

func (s *TCPTransfer) reconnect(clientKey string) {
	for {
		transferConn := s.getTransferConn(clientKey)
		if transferConn == nil {
			return
		}
		conn, err := net.Dial("tcp", s.dstAddress)
		if err == nil {
			s.mutex.Lock()
			transferConn.serverConn = conn
			s.clients[clientKey] = *transferConn
			s.mutex.Unlock()
			s.app.EventsEmit("transfer-tcp-info", clientKey, "connection reconnected")
			return
		}
		s.app.EventsEmit("transfer-tcp-info", clientKey, "trying to reconnect...")
		time.Sleep(RECONNECT_INTERVAL)
	}
}

func (s *TCPTransfer) handleEvents() {
	for {
		select {
		case <-s.shutdown:
			s.mutex.Lock()
			for addr, client := range s.clients {
				client.serverConn.Close()
				client.clientConn.Close()
				s.app.EventsEmit("transfer-tcp-info", addr, fmt.Sprintf("close connection %s", addr))
				fmt.Printf("Close connection %s\n", addr)
			}
			s.clients = make(map[string]TransferConn)
			s.mutex.Unlock()
			break
		case clientConn := <-s.addClient:
			if s.listener == nil {
				clientConn.Close()
				break
			}
			serverConn, err := net.Dial("tcp", s.dstAddress)
			if err != nil {
				s.app.EventsEmit("transfer-tcp-info", clientConn.RemoteAddr().String(), fmt.Sprintf("failed to connect to %s: %v", s.dstAddress, err))
				clientConn.Close()
				break
			}
			s.mutex.Lock()
			s.clients[clientConn.RemoteAddr().String()] = TransferConn{clientConn: clientConn, serverConn: serverConn}
			s.mutex.Unlock()
			go s.handleClientConnection(clientConn)
			go s.handleServerConnection(clientConn, serverConn)
		case conn := <-s.removeClient:
			clientKey := conn.RemoteAddr().String()
			s.mutex.Lock()
			transferConn, ok := s.clients[clientKey]
			delete(s.clients, clientKey)
			if ok {
				transferConn.serverConn.Close()
			}
			s.mutex.Unlock()
		case message := <-s.broadcastClient:
			s.mutex.RLock()
			for addr, client := range s.clients {
				_, err := client.clientConn.Write(message)
				if err != nil {
					s.app.EventsEmit("transfer-tcp-error", addr, fmt.Sprintf("error broadcasting message to client %s : %v", addr, err))
					fmt.Printf("Error broadcasting message to client %s: %s\n", addr, err.Error())
				}
			}
			s.mutex.RUnlock()
		case message := <-s.broadcastServer:
			s.mutex.RLock()
			for addr, client := range s.clients {
				_, err := client.serverConn.Write(message)
				if err != nil {
					s.app.EventsEmit("transfer-tcp-error", addr, fmt.Sprintf("error broadcasting message to client %s : %v", addr, err))
					fmt.Printf("Error broadcasting message to client %s: %s\n", addr, err.Error())
				}
			}
			s.mutex.RUnlock()
		}
	}
}

func (s *TCPTransfer) SendToServer(client string, message []byte) error {
	s.mutex.RLock()
	conn, ok := s.clients[client]
	s.mutex.RUnlock()
	if !ok {
		s.app.EventsEmit("transfer-tcp-error", client, fmt.Sprintf("client %s not found", client))
		return fmt.Errorf("client %s not found", client)
	}

	_, err := conn.serverConn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

func (s *TCPTransfer) SendToClient(client string, message []byte) error {
	s.mutex.RLock()
	conn, ok := s.clients[client]
	s.mutex.RUnlock()
	if !ok {
		s.app.EventsEmit("transfer-tcp-error", client, fmt.Sprintf("client %s not found", client))
		return fmt.Errorf("client %s not found", client)
	}

	_, err := conn.clientConn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

func (s *TCPTransfer) BroadcastToServer(message []byte) {
	s.broadcastServer <- message
}

func (s *TCPTransfer) BroadcastToClient(message []byte) {
	s.broadcastClient <- message
}
