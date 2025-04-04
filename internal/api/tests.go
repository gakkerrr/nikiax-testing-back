package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Test struct {
	ID            int    `json:"id,omitempty" db:"id"`
	Name          string `json:"name,omitempty" db:"name"`
	Alias         string `json:"alias,omitempty" db:"alias"`
	Result        int    `json:"result,omitempty" db:"result"`
	Err_txt       string `json:"err_txt,omitempty" db:"err_txt"`
	Set           int    `json:"set,omitempty" db:"set"`
	Time_duration int    `json:"time_duration,omitempty" db:"time_duration"`
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

		testId := chi.URLParam(r, "id")
		logger.Info("Получение всех теста " + testId + " из базы 	")

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

func CreateTests(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.Value("logger").(*slog.Logger)

		logger.Info("Добавление теста в базу данных")

		var tests []Test
		err := json.NewDecoder(r.Body).Decode(&tests)
		if err != nil {
			logger.Error("Ошибка получения данных из payload")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Info(fmt.Sprintf("Данные из payload добавления теста: %+v", tests))

		query := "SELECT MAX(id), MAX(set) FROM tests"
		row := db.QueryRow(query)

		var lastId int
		var lastSet int
		if err := row.Scan(&lastId, &lastSet); err != nil {
			logger.Error("Ошибка при получении последнего айди теста в ручке create_test", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Info("Последний айди в бд: " + strconv.Itoa(lastId))

		stmt := `INSERT INTO tests (id, name, alias, result, err_txt, set, time_duration) VALUES `
		values := []interface{}{}
		valuePlaceholders := []string{}

		for i, test := range tests {
			if test.Name == "" {
				logger.Error("Не найдено имя для теста в теле запроса", "index", i)
				http.Error(w, fmt.Sprintf("Не найдено имя для теста в теле запроса (индекс %d)", i), http.StatusBadRequest)
				return
			}

			nextId := lastId + i + 1
			nextSet := lastSet + 1
			curLen := len(values)

			valuePlaceholders = append(valuePlaceholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				curLen+1, curLen+2, curLen+3, curLen+4, curLen+5, curLen+6, curLen+7))

			values = append(values, nextId, test.Name, test.Alias, test.Result, test.Err_txt, nextSet, test.Time_duration)
		}

		stmt += strings.Join(valuePlaceholders, ", ")

		_, err = db.Exec(stmt, values...)
		if err != nil {
			logger.Error("Ошибка при добавлении тестов в бд в ручке create_test", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Все тесты успешно добавлены в базу данных")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(tests); err != nil {
			logger.Error("Ошибка преобразования данных в json формат в ручке create_test")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func DelTestId(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.Value("logger").(*slog.Logger)

		testId := chi.URLParam(r, "id")
		logger.Info("Удаление теста " + testId + " из базы 	")

		if testId == "" {
			logger.Error("Не получилось получить параметр для ручки tests/{id}")
			http.Error(w, http.StatusText(422), 422)
			return
		}

		query := "DELETE FROM tests WHERE id = " + testId
		rows, err := db.Query(query)
		if err != nil {
			logger.Error("Ошибка удаления данных из базы данных в ручке tests/id", "query", query)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
