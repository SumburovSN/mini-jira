package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mini-jira/task-service/internal/config"
	"mini-jira/task-service/internal/handler"
	"mini-jira/task-service/internal/middleware"
	"mini-jira/task-service/internal/repository"
	"mini-jira/task-service/internal/service"
)

func main() {
	cfg := config.Load()

	dbURL := cfg.DBUrl
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewTaskRepository(db)
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	r := chi.NewRouter()

	// Добавляем JWT Middleware для защиты маршрутов
	r.Use(middleware.JWTMiddleware([]byte(os.Getenv("JWT_SECRET"))))

	// Роуты
	r.Post("/tasks", h.Create)
	r.Get("/tasks", h.List)
	r.Delete("/tasks/{id}", h.Delete)
	r.Get("/tasks/{id}", h.GetById)
	r.Patch("/tasks/{id}", h.Update)

	port := cfg.Port
	if port == "" {
		port = "8083"
	}

	log.Println("Task service running on " + port)
	http.ListenAndServe(":"+port, r)
}
