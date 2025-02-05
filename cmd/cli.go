package cmd

import (
	"fmt"
	"os"
	"vovanDB/internal/database"
)

func Run() {
    // if len(os.Args) != 2 {
	// 	fmt.Println(database.ErrorArgs())
	// 	return
	// }

	// sql := os.Args[1]

	sql := "CREATE TABLE users (id int NOT NULL AUTO_INCREMENT,name text NULL,age int,is_admin bool,date date);"

	fmt.Println(database.Execute(sql))
}