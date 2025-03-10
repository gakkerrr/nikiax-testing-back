package tests

import (
	"database/sql"
	"fmt"
	"net/http"
)

func GetAllTests(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := db.Exec("SELECT * FROM tests")
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}
}
