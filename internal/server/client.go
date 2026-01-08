package server

import (
	"fmt"
	"net"
	"strings"
)

// ------------------------------------------------------------------|

type Client struct {
	conn net.Conn
	name string
}

// ------------------------------------------------------------------|

func newClient(c net.Conn, n string) Client {
	return Client{
		conn: c,
		name: n,
	}
}

// ------------------------------------------------------------------|

func client(conn net.Conn, addClient, deleteClient chan Client) {
	defer conn.Close()

	_, err := conn.Write([]byte("Enter your name: "))
	if err != nil {
		fmt.Printf("scan name: %v\n", err)
	}

	n := make([]byte, 4) // Надо придумать другой метод чтение
	_, err = conn.Read(n)
	if err != nil {
		fmt.Println(err)
		return
	}

	name := string(n)
	fmt.Printf("Name: %#v , len: %d\n", name, len(name))

	name = strings.TrimSpace(name)
	if len(name) == 0 {
		fmt.Println("The name is empty")
		return
	}

	client := newClient(conn, name)
	addClient <- client

	for {
		m := make([]byte, 1024)
		_, err := conn.Read([]byte(m))
		if err != nil {
			fmt.Printf("input error: %v\n", err)
			break
		}

		message := strings.TrimSpace(string(m))

		if len(message) == 0 {
			continue
		}
		fmt.Printf("INPUT: %s\n", message)
	}

	deleteClient <- client
}

// ------------------------------------------------------------------|
