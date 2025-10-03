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
	log.Debug("Подключение к серверу Postgres...")
	connFig := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
	conn, err := pgx.Connect(context.Background(), connFig)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к Postgres: %s", err)
	}
	return &Postgres{
		db:  conn,
		log: log,
	}, nil
}

// DbClose функция закрытия подключения к базе данных
func (pg *Postgres) DbClose() {
	pg.log.Debug("Отключение от сервера Postgres...")
	pg.db.Close(context.Background())
}
