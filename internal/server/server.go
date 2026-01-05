package server

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"tcp-chat/internal/models"
)

func RunServer() {
	listener, err := net.Listen("tcp", models.PORT)
	if err != nil {
		fmt.Printf("net.Listen: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("conn: %v\n", err)
			return
		}
		_, err = conn.Write([]byte("Hello, World!\n"))
		if err != nil {
			return
		}
		go client(conn)
	}
}

func client(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("_________YES________")

	for scanner.Scan() {
		fmt.Println("_________YES________")

		c, err := conn.Read([]byte(scanner.Text()))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(c)
	}
	fmt.Println("_________NO________")

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
}
