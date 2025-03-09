package main

import (
	"net/http"
	"nikiax-testing-back/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	db := database.MustLoadDB()
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Катя котик"))
	})
	http.ListenAndServe(":3000", r)
}
