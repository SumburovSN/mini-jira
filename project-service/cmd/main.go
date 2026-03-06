package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mini-jira/project-service/internal/config"
	"mini-jira/project-service/internal/handler"
	"mini-jira/project-service/internal/middleware"
	"mini-jira/project-service/internal/repository"
	"mini-jira/project-service/internal/service"
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

	repo := repository.NewProjectRepository(db)
	svc := service.NewProjectService(repo)
	h := handler.NewProjectHandler(svc)

	r := chi.NewRouter()

	// Добавляем JWT Middleware для защиты маршрутов
	r.Use(middleware.JWTMiddleware([]byte(os.Getenv("JWT_SECRET"))))

	// Роуты
	r.Post("/projects", h.Create)
	r.Get("/projects", h.List)
	r.Delete("/projects/{id}", h.Delete)
	r.Get("/projects/{id}", h.GetById)
	r.Patch("/projects/{id}", h.Update)

	port := cfg.Port
	if port == "" {
		port = "8082"
	}

	log.Println("Project service running on " + port)
	http.ListenAndServe(":"+port, r)
}
