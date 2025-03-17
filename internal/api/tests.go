package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Test struct {
	ID            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Alias         string `json:"alias" db:"alias"`
	Result        int    `json:"result" db:"result"`
	Err_txt       string `json:"err_txt" db:"err_txt"`
	Set           int    `json:"set" db:"set"`
	Time_duration int    `json:"time_duration" db:"time_duration"`
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
			if err := rows.Scan(&test.ID, &test.Name, &test.Alias, &test.Result, &test.Err_txt, &test.Set, &test.Time_duration); err != nil {
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

func GetTestId(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.Value("logger").(*slog.Logger)
		logger.Info("Получение всех тестов из базы данных")

		testId := chi.URLParam(r, "id")

		if testId == "" {
			logger.Error("Не получилось получить параметр для ручки tests/{id}")
			http.Error(w, http.StatusText(422), 422)
			return
		}

		query := "SELECT * FROM tests WHERE id = " + testId
		rows, err := db.Query(query)
		if err != nil {
			logger.Error("Ошибка получения данных из базы данных в ручке tests/id", "query", query)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tests []Test
		for rows.Next() {
			var test Test
			if err := rows.Scan(&test.ID, &test.Name, &test.Alias, &test.Result, &test.Err_txt, &test.Set, &test.Time_duration); err != nil {
				logger.Error("Ошибка при сканировании строки в ручке tests/id", "error", err)
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
			logger.Error("Ошибка преобразования данных в json формат в ручке tests/id")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
