package repository

import (
	"context"
	"errors"
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

func (r *UserRepository) UpdateEmail(ctx context.Context, userID int, email string) error {
	cmd, err := r.db.Exec(ctx,
		`UPDATE auth.users
		 SET email=$2
		 WHERE id=$1`,
		userID, email,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdatePasswordHash(ctx context.Context, userID int, hash string) error {
	cmd, err := r.db.Exec(ctx,
		`UPDATE auth.users
		 SET password_hash=$2
		 WHERE id=$1`,
		userID, hash,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

var ErrUserNotFound = errors.New("user not found")
