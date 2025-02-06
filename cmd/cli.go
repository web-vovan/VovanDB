package cmd

import (
	"fmt"
	"os"
	"vovanDB/internal/database"
)

func Run() {
    if len(os.Args) != 2 {
		fmt.Println(database.ErrorArgs())
		return
	}

	sql := os.Args[1]

	fmt.Println(database.Execute(sql))
}