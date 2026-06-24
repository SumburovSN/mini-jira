package repository

import (
	"context"
	"errors"
	"mini-jira/project-service/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemberRepository struct {
	db *pgxpool.Pool
}

func NewMemberRepository(db *pgxpool.Pool) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) Create(ctx context.Context, m *model.Member) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO project.project_members (project_id, user_id, role)
		 VALUES ($1,$2,$3)
		 RETURNING id`,
		m.ProjectID, m.UserID, m.Role,
	).Scan(&m.ID)
}

func (r *MemberRepository) GetAllByProject(ctx context.Context, projectID int) ([]model.Member, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, project_id, user_id, role
         FROM project.project_members
         WHERE project_id=$1`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.Member
	for rows.Next() {
		var m model.Member
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *MemberRepository) Delete(ctx context.Context, memberID int) error {
	cmd, err := r.db.Exec(ctx,
		`DELETE FROM project.project_members
         WHERE id=$1`,
		memberID,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrMemberNotFoundOrForbidden
	}

	return nil
}

func (r *MemberRepository) GetById(ctx context.Context, id int) (*model.Member, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, project_id, user_id, role
         FROM project.project_members
         WHERE id=$1`,
		id,
	)

	var m model.Member
	err := row.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	return &m, nil
}

func (r *MemberRepository) UpdateRole(ctx context.Context, memberID int, role model.Role) error {
	cmd, err := r.db.Exec(ctx,
		`UPDATE project.project_members
		 SET role=$2
		 WHERE id=$1`,
		memberID,
		role,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrMemberNotFoundOrForbidden
	}

	return nil
}

func (r *MemberRepository) GetUserRole(ctx context.Context, projectID int, userID int) (model.Role, error) {
	row := r.db.QueryRow(ctx,
		`SELECT role
		 FROM project.project_members
		 WHERE project_id=$1 AND user_id=$2`,
		projectID, userID,
	)

	var role model.Role

	err := row.Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrMemberNotInProject
		}
		return "", err
	}

	return role, nil
}

var ErrMemberNotFound = errors.New("member not found")

var ErrMemberNotFoundOrForbidden = errors.New("member not found or forbidden")

var ErrMemberNotInProject = errors.New("user is not a member of project")
