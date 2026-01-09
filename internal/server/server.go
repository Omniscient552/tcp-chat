package server

import (
	"fmt"
	"net"
	"time"

	"tcp-chat/internal/models"
)

// ------------------------------------------------------------------|

type Server struct {
	client       map[net.Conn]string // conn-name
	addClient    chan Client
	deleteClient chan Client
	broadcast    chan Message
}

// ------------------------------------------------------------------|

func NewServer() *Server {
	return &Server{
		client:       make(map[net.Conn]string),
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

		go client(conn, s.addClient, s.deleteClient, s.broadcast)
	}
}

// ------------------------------------------------------------------|

func (s *Server) manager() {
	for {
		select {
		case c := <-s.addClient:
			if len(s.client) >= models.MAX_CLIENT {
				sendMessage(c.conn, models.CHAT_FULL)
				c.conn.Close()
				continue
			}

			var isExists bool
			for _, name := range s.client { // Как-то так
				if c.name == name {
					isExists = true
					sendMessage(c.conn, models.NAME_TAKEN)
					break
				}
			}

			if isExists {
				go client(c.conn, s.addClient, s.deleteClient, s.broadcast)
				continue
			}

			s.client[c.conn] = c.name
			fmt.Println("len: ", len(s.client))
			fmt.Println("Connection client: ", c.name)

		case c := <-s.deleteClient:
			fmt.Println("Disconnect client: ", c.name)
			c.conn.Close()
			delete(s.client, c.conn)

		case newMessage := <-s.broadcast:
			for conn := range s.client {
				message := formatMessage(newMessage.name, newMessage.msg)
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
