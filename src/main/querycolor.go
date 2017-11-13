package main

import (
	"os"
	"fmt"
	"strconv"
	"log"
)

func main() {
	InitDB()
	defer CloseDB()

	amount, err := strconv.ParseFloat(os.Args[2], 32)
	if (err != nil) {
		log.Fatal("Usage: querycolor <color> <amount>")
	}
	res := QueryColor(os.Args[1], amount)
	for _, filename := range(res) {
		fmt.Println(filename)
	}
}
