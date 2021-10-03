package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

var (
	chatLog     []string
	connections []net.Conn
)

func server() {
	s, err := net.Listen("tcp", ":9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		c, err := s.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}
		connections = append(connections, c)
		go handleClient(c)
	}
}

func handleClient(c net.Conn) {
	var msg string
	for {
		err := gob.NewDecoder(c).Decode(&msg)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(msg)

		chatLog = append(chatLog, msg)

		for i := 0; i < len(connections); i++ {
			if c != connections[i] {
				err := gob.NewEncoder(connections[i]).Encode(msg)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

}

func backupMessages() {
	f, err := os.Create("Messages.txt")

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	for _, messages := range chatLog {
		f.WriteString(messages + "\n")
	}
}

func main() {
	var opc int

	fmt.Println("Starting server")
	go server()

	for {
		fmt.Println("\nMenu" +
			"\n 1. Backup messages" +
			"\n 2. Stop server")

		fmt.Scanln(&opc)

		switch opc {
		case 1:
			backupMessages()

		case 2:
			fmt.Println("\nGoodbye")
			return

		default:
			fmt.Println("\nWrong option")
		}
	}

}
