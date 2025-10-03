package main

import (
	"fmt"
	"junior_effectivemobile/db"
	"junior_effectivemobile/handlers"
	"junior_effectivemobile/server"
)

func main() {
	postgres, err := db.NewConPostgres()
	if err != nil {
		fmt.Print(err)
	}
	defer postgres.DbClose()
	handlers := handlers.NewHandlersApp(postgres)
	server := server.NewServer(handlers)
	server.ServerStart()
}
