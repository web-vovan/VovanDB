package cmd

import (
	"fmt"
	"os"
	"vovanDB/internal"
)

func Run() {
    if len(os.Args) != 2 {
		fmt.Println(internal.ErrorArgs())
		return
	}

	sql := os.Args[1]

	fmt.Println(internal.Execute(sql))
}