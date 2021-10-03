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

//global variables
var (
	chatLog        []string   //chat log keeps all messages
	userNames      []string   //userNames keep the clients names
	connections    []net.Conn //keeps every connection
	newConnections []net.Conn //auxiliary slice
)

func server() {
	//listen to client connection
	s, err := net.Listen("tcp", ":9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	//keeps on listening
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

		//if the message containg a "." means that it's a file
		if strings.Contains(msg, ".") {
			serverSendFile(c, msg)
		}

		//when the message contains that, means that it's its first connectio
		//we add the username to the slice
		if strings.Contains(msg, "has entered the chat room") {
			receive := strings.Split(msg, " has")
			userNames = append(userNames, receive[0])
		}

		//if it's not a disconnection we send the message to all clients
		//excluding the sender
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
		} else { //if the message is disconnect we close the connection
			//and we erase it from the slice
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
	//for every client except the sender we create a copy of the file sent
	receive := strings.Split(msg, ": ")
	dataBytes, err := ioutil.ReadFile(receive[1])
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(connections); i++ {
		if c != connections[i] {
			f, err := os.Create(userNames[i] + "/" + receive[1])
			if err != nil {
				fmt.Println(err)
			}

			defer f.Close()

			f.Write(dataBytes)
		}
	}

	//same for the server
	f, err := os.Create("server/" + receive[1])
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	f.Write(dataBytes)
}

//when the backup it's needed we do it in the reserved
//server directory
func backupMessages() {
	f, err := os.Create("server/Messages.txt")

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

	//we create the server directory
	err := os.Mkdir("server/", 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	//goroutine server
	go server()

	//keeps on listening server-client input
	for {
		fmt.Println("\nMenu" +
			"\n 1. Backup messages" +
			"\n 2. Stop server")

		fmt.Scanln(&opc)

		//backup chat log
		switch opc {
		case 1:
			backupMessages()

		//goodbye message
		case 2:
			fmt.Println("\nGoodbye")
			return

		//any other input won't be accepted
		default:
			fmt.Println("\nWrong option")
		}
	}

}
