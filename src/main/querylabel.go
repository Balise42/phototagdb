package main

import (
	"os"
	"fmt"
)

func main() {
	InitDB()
	defer CloseDB()

	res := QueryLabels(os.Args[1:])
	for _, filename := range(res) {
		fmt.Println(filename)
	}
}
