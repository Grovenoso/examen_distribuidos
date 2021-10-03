package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var (
	chatLog        []string
	userNames      []string
	connections    []net.Conn
	newConnections []net.Conn
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

		if strings.Contains(msg, ".") {
			serverSendFile(c, msg)
		}

		if strings.Contains(msg, "has entered the room") {
			receive := strings.Split(msg, " has")
			userNames = append(userNames, receive[0])
		}

		if msg != "disconnect" {
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
		} else {
			for i := 0; i < len(connections); i++ {
				if c == connections[i] {
					c.Close()
					fmt.Println("Client has disconnected")
				} else {
					newConnections = append(newConnections, connections[i])
				}
			}
			connections = newConnections
			newConnections = nil
		}
	}

}

func serverSendFile(c net.Conn, msg string) {
	for i := 0; i < len(connections); i++ {
		if c != connections[i] {
			receive := strings.Split(msg, ": ")
			dataBytes, err := ioutil.ReadFile(receive[1])
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.Create(userNames[i] + "/" + receive[1])
			if err != nil {
				fmt.Println(err)
			}

			defer f.Close()

			f.Write(dataBytes)
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
