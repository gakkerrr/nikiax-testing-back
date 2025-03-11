package tests

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Test struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Alias  string `json:"alias" db:"alias"`
	Result int    `json:"result" db:"result"`
	Stage  int    `json:"stage" db:"stage"`
}

func GetAllTests(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM tests")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tests []Test
		for rows.Next() {
			var test Test
			if err := rows.Scan(&test.ID, &test.Name, &test.Alias, &test.Result, &test.Stage); err != nil {
				log.Printf("Ошибка при сканировании строки: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Сканированная строка: %+v", test)
			tests = append(tests, test)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(tests); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
