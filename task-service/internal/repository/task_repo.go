package repository

import (
	"context"
	"errors"
	"mini-jira/task-service/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, t *model.Task) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO task.tasks (project_id, title, description, status, assignee_id)
         VALUES ($1,$2,$3,$4,$5)
         RETURNING id, created_at`,
		t.ProjectID, t.Title, t.Description, t.Status, t.AssigneeID,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *TaskRepository) GetAllByAssignee(ctx context.Context, assigneeID int) ([]model.Task, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, project_id, title, description, status, assignee_id, created_at
         FROM task.tasks
         WHERE assignee_id=$1`,
		assigneeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.AssigneeID, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int, assigneeID int) error {
	cmd, err := r.db.Exec(ctx,
		`DELETE FROM task.tasks
         WHERE id=$1 AND assignee_id=$2`,
		id, assigneeID,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrProjectNotFoundOrForbidden
	}

	return nil
}

var ErrProjectNotFound = errors.New("task not found")

// var ProjectNameEmpty = errors.New("project name is empty")
var ErrProjectNotFoundOrForbidden = errors.New("task not found or forbidden")

func (r *TaskRepository) GetById(ctx context.Context, id int, assigneeID int) (*model.Task, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, project_id, title, description, status, assignee_id, created_at
         FROM task.tasks
         WHERE id=$1 AND assignee_id=$2`,
		id, assigneeID,
	)

	var t model.Task
	err := row.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.AssigneeID, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return &t, nil
}

func (r *TaskRepository) Update(ctx context.Context, title string, description string, status string, id int, assigneeID int) error {
	cmd, err := r.db.Exec(ctx,
		`UPDATE task.tasks
		 SET title=$1, description=$2, status=$3 
         WHERE id=$4 AND assignee_id=$5`,
		title, description, status, id, assigneeID,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrProjectNotFoundOrForbidden
	}

	return nil
}
