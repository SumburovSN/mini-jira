package service

import (
	"context"
	"time"

	"mini-jira/auth-service/internal/model"
	"mini-jira/auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(r *repository.UserRepository, secret string) *AuthService {
	return &AuthService{repo: r, jwtSecret: []byte(secret)}
}

func (s *AuthService) Register(ctx context.Context, email, password string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return s.repo.Create(ctx, email, string(hash))
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}
