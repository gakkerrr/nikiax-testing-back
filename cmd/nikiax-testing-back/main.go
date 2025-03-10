package main

import (
	"net/http"
	tests "nikiax-testing-back/internal/api"
	"nikiax-testing-back/internal/database"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	db := database.MustLoadDB()
	defer db.Close()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Катя где фронт"))
	})

	r.Get("/tests", tests.GetAllTests(db))

	http.ListenAndServe(":3000", r)
}
