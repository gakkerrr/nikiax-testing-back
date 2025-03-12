package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Test struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Alias  string `json:"alias" db:"alias"`
	Result int    `json:"result" db:"result"`
	Stage  int    `json:"stage" db:"stage"`
}

func GetAllTests(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.Value("logger").(*slog.Logger)
		logger.Info("Получение всех тестов из базы данных")

		query := "SELECT * FROM tests"
		rows, err := db.Query(query)
		if err != nil {
			logger.Error("Ошибка получения данных из базы данных в ручке tests", "query", query)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tests []Test
		for rows.Next() {
			var test Test
			if err := rows.Scan(&test.ID, &test.Name, &test.Alias, &test.Result, &test.Stage); err != nil {
				logger.Error("Ошибка при сканировании строки в ручке tests", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tests = append(tests, test)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(tests); err != nil {
			logger.Error("Ошибка преобразования данных в json формат в ручке tests")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
