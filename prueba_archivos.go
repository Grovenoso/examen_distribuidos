package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	dataBytes, err := ioutil.ReadFile("2049.jpg")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir("user/", 0755)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create("user/2049_1.jpg")

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	f.Write(dataBytes)
}
