package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"mini-jira/auth-service/internal/config"
	"mini-jira/auth-service/internal/handler"
	"mini-jira/auth-service/internal/repository"
	"mini-jira/auth-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	db, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo, cfg.JWTSecret)
	h := handler.NewAuthHandler(svc)

	r := chi.NewRouter()
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Group(func(r chi.Router) {
		r.Use(handler.JWTMiddleware([]byte(cfg.JWTSecret)))
		r.Get("/me", h.Me)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Println("Auth service running on", port)
	http.ListenAndServe(":"+port, r)
}
