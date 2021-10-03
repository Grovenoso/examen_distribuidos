package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

func runClient(userName string) {
	c, err := net.Dial("tcp", ":9999")

	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	defer c.Close()

	var message, received string

	//Welcome message
	message = userName + " has entered the room"
	err = gob.NewEncoder(c).Encode(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		fmt.Println("Ingrese su mensaje")
		fmt.Scanln(&message)
		message = userName + ": " + message

		err = gob.NewEncoder(c).Encode(message)
		if err != nil {
			fmt.Println(err)
			return
		}

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
	}
}

func main() {
	//var opc int
	var userName string

	fmt.Println("Enter your username")
	fmt.Scanln(&userName)

	go runClient(userName)

	/*
		for {
			fmt.Println("\nMenu" +
				"\n 1. Send message" +
				"\n 2. Send File" +
				"\n 3. Stop client")

			fmt.Scanln(&opc)

			switch opc {
			case 1:
				//send Message

			case 2:
				//send File

			case 3:
				fmt.Println("\nGoodbye")
				return

			default:
				fmt.Println("\nWrong option")
			}
		}
	*/
	for {
	}
}
