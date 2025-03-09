package main

import (
	"net/http"
	"nikiax-testing-back/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"
)

func main() {

	db := database.MustLoadDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO tests (name, alias, result, stage) VALUES ($1, $2, $3, $4)", "t5", "test t5", 0, 5)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Катя котик"))
	})
	http.ListenAndServe(":3000", r)
}
