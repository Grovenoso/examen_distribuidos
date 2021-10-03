package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func runClient(userName string, status chan int, msg chan string) {
	c, err := net.Dial("tcp", ":9999")

	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	defer c.Close()

	var message, received string

	//welcome message
	message = userName + " has entered the room"

	//always listening
	go func() {
		for {
			err := gob.NewDecoder(c).Decode(&received)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(received)
		}
	}()

	//keep the connection
	for {
		select {
		case _status := <-status:
			//first connection
			if _status == 0 {
				err = gob.NewEncoder(c).Encode(message)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			//send message
			if _status == 1 {
				message = <-msg
				err := gob.NewEncoder(c).Encode(message)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			//send file
			if _status == 2 {
				sendFile(c, userName)
			}
		}
	}
}

func sendMessage(c net.Conn, userName string) {
	fmt.Println("Enter your message")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message := scanner.Text()

	message = userName + ": " + message

	err := gob.NewEncoder(c).Encode(message)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func sendFile(c net.Conn, userName string) {
	fmt.Println("Not finished yet")
}

func main() {
	var opc int
	var status = make(chan int)
	var msg = make(chan string)

	fmt.Println("Enter your username")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userName := scanner.Text()

	go runClient(userName, status, msg)
	status <- 0

	fmt.Println("\nMenu" +
		"\n 1. Send message" +
		"\n 2. Send File" +
		"\n 3. Stop client")

	for {
		fmt.Scanln(&opc)

		switch opc {
		case 1:
			fmt.Println("Enter your message")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			message := scanner.Text()
			message = userName + ": " + message
			status <- 1
			msg <- message

		case 2:
			//status <- 2

		case 3:
			fmt.Println("\nGoodbye")
			return

		default:
			fmt.Println("\nWrong option")
		}
	}
}
