package database

import (
	"context"
	"fmt"
	"log/slog"

	"database/sql"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
)

type dbConfig struct {
	Port     int    `env:"DB_PORT" env-default:"5432"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	DbName   string `env:"DB_NAME" env-default:"postgres"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD"`
}

func MustLoadDB(ctx context.Context) *sql.DB {

	logger := ctx.Value("logger").(*slog.Logger)

	var cfg dbConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		logger.Error("Ошибка считывания переменных окружения", "error", err)
		panic(err)
	}

	logger.Info("Данные для базы данных успешно считаны из переменных среды")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
	)

	db, err := sql.Open(cfg.DbName, psqlInfo)
	if err != nil {
		logger.Error("Ошибка при подключении к базе данных", "error", err)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		logger.Error("Ошибка при пинге базы данных", "error", err)
		panic(err)
	}

	return db
}

// TODO: сюда добавить функции по получению и добавлению данных в бд, т.е. убрать всю логику из tests.go сюда. Там оставить только вызов и обработку даннных
