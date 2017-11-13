package main

import (
	"os"
	"fmt"
)

func main() {
	InitDB()
	defer CloseDB()

	res := QueryText(os.Args[1])
	for _, filename := range(res) {
		fmt.Println(filename)
	}
}
