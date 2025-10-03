package db

import (
	"context"
	"fmt"
	"junior_effectivemobile/config"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// Postgres структура для наследования методов
type Postgres struct {
	db  *pgx.Conn
	log *logrus.Logger
}

// NewConPostgres функция подключения к базе данных, возвращает структуру Postgres
func NewConPostgres(log *logrus.Logger, cfg *config.Config) (*Postgres, error) {
	log.Debug("Начало подключения к серверу Postgres")

	connFig := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
	log.Debug("Формирование строки подключения",
		"host", cfg.DbHost,
		"port", cfg.DbPort,
		"database", cfg.DbName,
		"user", cfg.DbUser)

	conn, err := pgx.Connect(context.Background(), connFig)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к Postgres: %s", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения с Postgres: %s", err)
	}

	log.Debug("Подключение к Postgres успешно установлено")
	return &Postgres{
		db:  conn,
		log: log,
	}, nil
}

// DbClose функция закрытия подключения к базе данных
func (pg *Postgres) DbClose() {
	pg.log.Debug("Закрытие соединения с базой данных")
	err := pg.db.Close(context.Background())
	if err != nil {
		pg.log.Debug("Ошибка при закрытии соединения", "error", err)
	} else {
		pg.log.Debug("Отключение от сервера Postgres успешно завершено")
	}
}
