package repository

import (
	"context"
	"errors"
	"mini-jira/project-service/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, p *model.Project) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO project.projects (name, owner_id)
         VALUES ($1,$2)
         RETURNING id, created_at`,
		p.Name, p.OwnerID,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *ProjectRepository) GetAllByOwner(ctx context.Context, ownerID int) ([]model.Project, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, owner_id, created_at
         FROM project.projects
         WHERE owner_id=$1`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.OwnerID, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id int, ownerID int) error {
	cmd, err := r.db.Exec(ctx,
		`DELETE FROM project.projects
         WHERE id=$1 AND owner_id=$2`,
		id, ownerID,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrProjectNotFoundOrForbidden
	}

	return nil
}

var ErrProjectNotFound = errors.New("project not found")

// var ProjectNameEmpty = errors.New("project name is empty")
var ErrProjectNotFoundOrForbidden = errors.New("project not found or forbidden")

func (r *ProjectRepository) GetById(ctx context.Context, id int, ownerID int) (*model.Project, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, owner_id, created_at
         FROM project.projects
         WHERE id=$1 AND owner_id=$2`,
		id, ownerID,
	)

	var p model.Project
	err := row.Scan(&p.ID, &p.Name, &p.OwnerID, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepository) Update(ctx context.Context, name string, id int, ownerID int) error {
	// if name == "" {
	// 	return ProjectNameEmpty
	// }
	cmd, err := r.db.Exec(ctx,
		`UPDATE project.projects
		 SET name=$1 
         WHERE id=$2 AND owner_id=$3`,
		name, id, ownerID,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrProjectNotFoundOrForbidden
	}

	return nil
}
