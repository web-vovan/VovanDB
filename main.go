package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(ErrorArgs())
		return
	}

	sql := os.Args[1]

	fmt.Println(Execute(sql))
}
