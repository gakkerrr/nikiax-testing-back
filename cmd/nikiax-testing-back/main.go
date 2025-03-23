package main

import (
	"context"
	"log/slog"
	"net/http"
	tests "nikiax-testing-back/internal/api"
	"nikiax-testing-back/internal/database"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", logger)

	db := database.MustLoadDB(ctx)
	logger.Info("База данных успешно подключена")
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

	r.Get("/tests", tests.GetAllTests(ctx, db))
	r.Get("/tests/{id}", tests.GetTestId(ctx, db))
	r.Post("/create_tests", tests.CreateTests(ctx, db))

	logger.Info("Сервер успешно запущен")
	http.ListenAndServe(":3000", r)
}
