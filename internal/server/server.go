package server

import (
	"fmt"
	"net"

	"tcp-chat/internal/models"
)

// ------------------------------------------------------------------|

type Server struct {
	client       map[net.Conn]string // conn-name
	addClient    chan Client
	deleteClient chan Client
	broadcast    chan string
}

// ------------------------------------------------------------------|

func NewServer() *Server {
	return &Server{
		client:       make(map[net.Conn]string),
		addClient:    make(chan Client, 4),
		deleteClient: make(chan Client, 4),
		broadcast:    make(chan string, 4),
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

	fmt.Println("Server is listening...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("conn: %v\n", err)
			return
		}

		go client(conn, s.addClient, s.deleteClient)
	}
}

// ------------------------------------------------------------------|

func (s *Server) manager() {
	for {
		select {
		case c := <-s.addClient: // подключение клиента
			fmt.Println("Connection client: ", c.name)
			s.client[c.conn] = c.name
			fmt.Println("All client: ", len(s.client))

		case c := <-s.deleteClient: // отключение клиента
			fmt.Println("Disconnect client: ", c.name)
			delete(s.client, c.conn)
			fmt.Println("All client: ", len(s.client))

		case m := <-s.broadcast: // отправка всем сообщение
			fmt.Println("Message: ", m)

			for c := range s.client {
				sendMessage(c, m)
			}
		}
	}
}

// ------------------------------------------------------------------|

func sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error send message: %v\n", err)
	}
}

// ------------------------------------------------------------------|
