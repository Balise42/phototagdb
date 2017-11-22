package main

import (
	"fmt"
	"os"
)

func main() {
	InitDB()
	defer CloseDB()

	res := QueryImage(os.Args[1])
	for _, tag := range res {
		fmt.Println(tag)
	}
}
