package mircat

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type TCPServer struct {
	address      string
	listener     net.Listener
	clients      map[string]net.Conn
	mutex        sync.RWMutex
	broadcast    chan []byte
	addClient    chan net.Conn
	removeClient chan net.Conn
	shutdown     chan bool
	app          *App
}

func NewTCPServer(app *App) *TCPServer {
	return &TCPServer{
		clients:      make(map[string]net.Conn),
		broadcast:    make(chan []byte),
		addClient:    make(chan net.Conn),
		removeClient: make(chan net.Conn),
		shutdown:     make(chan bool),
		app:          app,
	}
}

func (s *TCPServer) Start(address string) error {
	var err error
	s.address = address
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	fmt.Printf("Listening on %s\n", s.address)

	go s.handleEvents()

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					s.app.EventsEmit("server-tcp-error", "server", fmt.Sprintf("error accepting connection: %v", err))
				}
				fmt.Printf("Error accepting connection: %s\n", err.Error())
				return
			}
			s.app.EventsEmit("server-tcp-info", conn.RemoteAddr().String(), fmt.Sprintf("client connected: %s", conn.RemoteAddr()))
			fmt.Printf("New client connected: %s\n", conn.RemoteAddr())

			s.addClient <- conn
		}
	}()
	return nil
}

func (s *TCPServer) Stop() {
	if s.listener != nil {
		s.shutdown <- true

		s.listener.Close()
		s.listener = nil
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.app.EventsEmit("server-tcp-info", conn.RemoteAddr().String(), fmt.Sprintf("client disconnected: %s", conn.RemoteAddr()))
		fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())

		s.removeClient <- conn
	}()

	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF && !errors.Is(err, net.ErrClosed) {
				s.app.EventsEmit("server-tcp-error", conn.RemoteAddr().String(), fmt.Sprintf("error reading from client %s : %v", conn.RemoteAddr(), err))
			}
			fmt.Printf("Error reading from client %s: %s\n", conn.RemoteAddr(), err.Error())
			return
		}
		message := append([]byte{}, buffer[:n]...)
		s.app.EventsEmit("server-tcp-data", conn.RemoteAddr().String(), message)
		//s.broadcast <- message
	}
}

func (s *TCPServer) handleEvents() {
	for {
		select {
		case <-s.shutdown:
			s.mutex.Lock()
			for addr, client := range s.clients {
				client.Close()
				s.app.EventsEmit("server-tcp-info", addr, fmt.Sprintf("close connection %s", addr))
				fmt.Printf("Close connection %s\n", addr)
			}
			s.clients = make(map[string]net.Conn)
			s.mutex.Unlock()
			break
		case conn := <-s.addClient:
			if s.listener == nil {
				conn.Close()
				break
			}
			s.mutex.Lock()
			s.clients[conn.RemoteAddr().String()] = conn
			s.mutex.Unlock()
			go s.handleConnection(conn)
		case conn := <-s.removeClient:
			s.mutex.Lock()
			delete(s.clients, conn.RemoteAddr().String())
			s.mutex.Unlock()
		case message := <-s.broadcast:
			s.mutex.RLock()
			for addr, client := range s.clients {
				_, err := client.Write(message)
				if err != nil {
					s.app.EventsEmit("server-tcp-error", addr, fmt.Sprintf("error broadcasting message to client %s : %v", addr, err))
					fmt.Printf("Error broadcasting message to client %s: %s\n", addr, err.Error())
				}
			}
			s.mutex.RUnlock()
		}
	}
}

func (s *TCPServer) SendMessage(client string, message []byte) error {
	s.mutex.RLock()
	conn, ok := s.clients[client]
	s.mutex.RUnlock()
	if !ok {
		s.app.EventsEmit("server-tcp-error", client, fmt.Sprintf("client %s not found", client))
		return fmt.Errorf("client %s not found", client)
	}

	_, err := conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

func (s *TCPServer) BroadcastMessage(message []byte) {
	s.broadcast <- message
}
