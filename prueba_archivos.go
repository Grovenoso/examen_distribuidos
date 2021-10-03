package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	dataBytes, err := ioutil.ReadFile("2049.jpg")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		n := strconv.Itoa(i)
		f, err := os.Create("2049_" + n + ".jpg")

		if err != nil {
			fmt.Println(err)
		}

		defer f.Close()

		f.Write(dataBytes)
	}
}
