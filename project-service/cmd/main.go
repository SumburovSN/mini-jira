package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"mini-jira/project-service/internal/config"
	"mini-jira/project-service/internal/handler"
	"mini-jira/project-service/internal/middleware"
	"mini-jira/project-service/internal/repository"
	"mini-jira/project-service/internal/router"
	"mini-jira/project-service/internal/service"
)

func main() {
	cfg := config.Load()

	dbURL := cfg.DBUrl
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	jwtSecret := cfg.JWTSecret
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	memberRepo := repository.NewMemberRepository(db)
	memberService := service.NewMemberService(memberRepo)
	memberHandler := handler.NewMemberHandler(memberService)

	r := chi.NewRouter()

	// Добавляем JWT Middleware для защиты маршрутов
	r.Use(middleware.JWTMiddleware([]byte(jwtSecret)))

	router.RegisterProjectRoutes(r, projectHandler)
	router.RegisterMemberRoutes(r, memberHandler)

	port := cfg.Port
	if port == "" {
		port = "8082"
	}

	log.Println("Project service running on " + port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
