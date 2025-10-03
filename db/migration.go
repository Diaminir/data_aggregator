package db

import (
	"context"
	"fmt"
)

// RunMigration функция запускающая миграцию
func (pg *Postgres) RunMigration() error {
	pg.log.Debug("Начало миграции базы данных")

	sql := `CREATE TABLE IF NOT EXISTS subscriptions
(
	id serial PRIMARY KEY,
	service_name varchar(255) NOT NULL,
    price int NOT NULL CHECK (price > 0),
    user_id uuid NOT NULL,
    start_date date NOT NULL,
    end_date date NULL
)`

	if _, err := pg.db.Exec(context.Background(), sql); err != nil {
		return fmt.Errorf("ошибка миграции базы данных: %s", err)
	}

	pg.log.Debug("Миграция успешно завершена")
	return nil
}
