package database

import (
	"fmt"

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

func MustLoadDB() *sql.DB {

	var cfg dbConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
	)

	fmt.Println(psqlInfo)

	db, err := sql.Open(cfg.DbName, psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
