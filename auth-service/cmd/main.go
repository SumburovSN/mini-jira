package main

import (
	"context"
	"log"
	"net/http"

	"mini-jira/auth-service/internal/config"
	"mini-jira/auth-service/internal/handler"
	"mini-jira/auth-service/internal/middleware"
	"mini-jira/auth-service/internal/repository"
	"mini-jira/auth-service/internal/router"
	"mini-jira/auth-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo, jwtSecret)
	h := handler.NewAuthHandler(svc)

	r := chi.NewRouter()
	// Публичные маршруты
	router.RegisterPublicAuthRoutes(r, h)

	// Защищённые маршруты
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware([]byte(jwtSecret)))
		router.RegisterPrivateAuthRoutes(r, h)
	})

	port := cfg.Port
	if port == "" {
		port = "8081"
	}

	log.Println("Auth service running on", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
