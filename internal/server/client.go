package server

import (
	"fmt"
	"io"
	"net"
	"strings"

	"tcp-chat/internal/models"
)

// ------------------------------------------------------------------|

type Client struct {
	conn   net.Conn
	name   string
	change bool
}

// ------------------------------------------------------------------|

type Message struct {
	name string
	msg  string
}

// ------------------------------------------------------------------|

func newClient(c net.Conn, n string) Client {
	return Client{
		conn: c,
		name: n,
	}
}

// ------------------------------------------------------------------|

func client(conn net.Conn, addClient, deleteClient chan<- Client, broadcast chan<- Message) {
	client, err := createClient(conn)
	if err != nil {
		fmt.Println(err)
	}

	addClient <- client

	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("read error: %v\n", err)
			}
			break
		}

		message := strings.TrimSpace(string(b[:n]))

		if len(message) == 0 {
			continue
		}

		if strings.HasPrefix(message, "/name") {
			part := strings.Split(message, " ")
			newName := strings.TrimSpace(part[1])

			if newName == "" {
				sendMessage(client.conn, models.EMPTY_NAME)
				continue
			}
			client.name = newName
			client.change = true

			addClient <- client
			continue
		}

		newMessage := Message{
			name: client.name,
			msg:  message,
		}

		broadcast <- newMessage
	}

	deleteClient <- client
}

// ------------------------------------------------------------------|

func createClient(conn net.Conn) (Client, error) {
	_, err := conn.Write([]byte("[ENTER YOUR NAME]: "))
	if err != nil {
		return Client{}, fmt.Errorf("write error: %v", err)
	}

	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		return Client{}, fmt.Errorf("read error: %v", err)
	}

	name := strings.TrimSpace(string(b[:n]))

	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return Client{}, fmt.Errorf("the name is empty")
	}

	return newClient(conn, name), nil
}

// ------------------------------------------------------------------|
