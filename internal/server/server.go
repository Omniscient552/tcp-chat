package server

import (
	"fmt"
	"net"
	"time"

	"tcp-chat/internal/models"
)

// ------------------------------------------------------------------|

type Server struct {
	client       map[string]chan string // name-chan
	addClient    chan Client
	deleteClient chan Client
	broadcast    chan Message
}

// ------------------------------------------------------------------|

func NewServer() *Server {
	return &Server{
		client:       make(map[string]chan string),
		addClient:    make(chan Client, 4),
		deleteClient: make(chan Client, 4),
		broadcast:    make(chan Message, 4),
	}
}

// ------------------------------------------------------------------|

func RunServer() {
	listener, err := net.Listen("tcp", models.PORT)
	if err != nil {
		fmt.Printf("net.Listen: %v\n", err)
		return
	}
	defer listener.Close()

	s := NewServer()
	go s.manager()

	fmt.Printf("The server is listening on port %s...\n", models.PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("conn: %v\n", err)
			return
		}

		go reader(conn, s.addClient, s.deleteClient, s.broadcast)
	}
}

// ------------------------------------------------------------------|

func (s *Server) manager() {
	for {
		select {
		case c := <-s.addClient:
			if c.change {
				changeClientName(s.client, c.name, c.writeCh)
			}

			ok, msg := saveClient(s.client, c.name, c.writeCh)

			if !ok {
				sendMessage(c.conn, msg)
				c.conn.Close()
				continue
			}

			notification(s.client, msg)
			fmt.Println("Connection client: ", c.name)

		case c := <-s.deleteClient:
			c.conn.Close()
			delete(s.client, c.name)
			fmt.Println("Disconnect client: ", c.name)

			msg := c.name + models.UserLeft
			notification(s.client, msg)

		case newMessage := <-s.broadcast:
			message := formatMessage(newMessage.name, newMessage.msg)
			for conn := range s.client {
				sendMessage(conn, message)
			}

		}
	}
}

// ------------------------------------------------------------------|

func formatMessage(name, message string) string {
	t := time.Now().Format(time.DateTime)
	return fmt.Sprintf("\x1b[33;3;7m[%v]\x1b[32m[%s]:\x1b[34m %s\x1b[0m\n", t, name, message)
}

// ------------------------------------------------------------------|

func sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error send message: %v\n", err)
	}
}

// ------------------------------------------------------------------|

func notification(client map[string]chan string, msg string) {
	for _, ch := range client {
		// c <- msg
		select {
		case ch <- msg:
			continue
		default:
			go func() {
				select {
				case <-time.After(5 * time.Second):
				case ch <- msg:
				}
			}()
		}
	}
}

// ------------------------------------------------------------------|

func saveClient(client map[string]chan string, name string, writeCh chan string) (bool, string) {
	if len(client) >= models.MaxClinet {
		return false, models.ChatFull
	}

	if _, exists := client[name]; exists {
		return false, models.NameTaken
	}

	client[name] = writeCh

	msg := name + models.UserJoined
	return true, msg
}

// ------------------------------------------------------------------|

func changeClientName(client map[string]chan string, name string, writeCh chan string) {
	var oldName string
	for oldName, ch := range client {
		if ch == writeCh {
			delete(client, oldName)
			client[name] = writeCh
			break
		}
	}

	msg := fmt.Sprintf("User %s changed name to %s", oldName, name)

	go notification(client, msg)
}

// ------------------------------------------------------------------|
