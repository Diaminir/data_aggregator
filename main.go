package main

import (
	"context"
	"fmt"
	"junior_effectivemobile/config"
	"junior_effectivemobile/db"
	"junior_effectivemobile/handlers"
	"junior_effectivemobile/logger"
	"junior_effectivemobile/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log, file := logger.NewLog()
	defer file.Close()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	postgres, err := db.NewConPostgres(log, cfg)
	if err != nil {
		fmt.Print(err)
	}

	defer postgres.DbClose()
	handlers := handlers.NewHandlersApp(postgres, log)
	server := server.NewServer(handlers, log)
	srv := server.ServerStart()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Warn("Остановка сервера...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Принудительная остановка:", err)
	}
	log.Warn("Сервер завершил работу")
}
