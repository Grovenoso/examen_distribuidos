package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

//dialing the server
func runClient(userName string, status chan int, msg chan string) {
	c, err := net.Dial("tcp", ":9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	var message, received string

	//welcome message
	message = userName + " has entered the chat room"

	//create its own directory for files
	err = os.Mkdir(userName+"/", 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	//always listening
	go func() {
		defer c.Close()
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
				message = <-msg
				clientSendFile(c, userName, message)
			}

			//terminate connection
			if _status == 3 {
				message = <-msg
				err := gob.NewEncoder(c).Encode(message)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("\nGoodbye!")
				return
			}
		}
	}
}

//send file name and its extension
func clientSendFile(c net.Conn, userName, message string) {
	message = userName + ": " + message
	err := gob.NewEncoder(c).Encode(message)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	var opc int
	var status = make(chan int)
	var msg = make(chan string)

	//until the client enters their name
	//the client won't connect to the server
	fmt.Println("Enter your username")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userName := scanner.Text()

	//goroutine connecting to server
	go runClient(userName, status, msg)
	status <- 0

	//main menu
	fmt.Println("\nMenu" +
		"\n 1. Send message" +
		"\n 2. Send File" +
		"\n 3. Stop client")

	//keep listening to user input
	for {
		fmt.Scanln(&opc)

		switch opc {
		//send message
		case 1:
			fmt.Println("Enter your message")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			message := scanner.Text()
			message = userName + ": " + message
			status <- 1
			msg <- message
		//send file
		case 2:
			fmt.Println("Enter your file name")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			message := scanner.Text()
			status <- 2
			msg <- message
		//disconnect
		case 3:
			status <- 3
			msg <- "disconnect"
			return
		//any other input won't be accepted
		default:
			fmt.Println("\nWrong option")
		}
	}
}
