package main

import (
	"fmt"
	"os"
)

func main() {
	InitDB()
	defer CloseDB()

	res := QueryText(os.Args[1])
	for _, filename := range res {
		fmt.Println(filename)
	}
}
