package main

import (
	"context"
	"junior_effectivemobile/config"
	"junior_effectivemobile/db"
	_ "junior_effectivemobile/docs"
	"junior_effectivemobile/handlers"
	"junior_effectivemobile/logger"
	"junior_effectivemobile/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Aгрегатор данных
// @version 1.0
// @description REST-cервис для агрегации данных об онлайн-подписках пользователей
// @contact.email dima.drozdov15@mail.ru
// @host localhost:8080
// @BasePath /
func main() {
	log, file := logger.NewLog()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	postgres, err := db.NewConPostgres(log, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.DbClose()
	if err := postgres.RunMigration(); err != nil {
		log.Fatal(err)
	}
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
