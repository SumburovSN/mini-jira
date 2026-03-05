package repository

import (
	"context"
	"mini-jira/auth-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email, hash string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO auth.users(email, password_hash) VALUES($1,$2)`,
		email, hash)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash FROM auth.users WHERE email=$1`,
		email)

	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash FROM auth.users WHERE id=$1`, id)

	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
