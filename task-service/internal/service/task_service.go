package service

import (
	"context"
	"errors"
	"mini-jira/task-service/internal/model"
	"mini-jira/task-service/internal/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(r *repository.TaskRepository) *TaskService {
	return &TaskService{repo: r}
}

func (s *TaskService) Create(ctx context.Context, projectID int, title string, description string, status string, assigneeID int) (*model.Task, error) {
	t := &model.Task{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      status,
		AssigneeID:  assigneeID,
	}
	err := s.repo.Create(ctx, t)
	return t, err
}

func (s *TaskService) List(ctx context.Context, assigneeID int) ([]model.Task, error) {
	return s.repo.GetAllByAssignee(ctx, assigneeID)
}

func (s *TaskService) Delete(ctx context.Context, id int, assigneeID int) error {
	return s.repo.Delete(ctx, id, assigneeID)
}

var ErrTaskNotFound = errors.New("task not found")
var ErrTaskTitleEmpty = errors.New("task title is empty")
var ErrPTaskNotFoundOrForbidden = errors.New("task not found or forbidden")

func (s *TaskService) GetById(ctx context.Context, id int, assigneeID int) (*model.Task, error) {
	p, err := s.repo.GetById(ctx, id, assigneeID)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *TaskService) Update(ctx context.Context, title string, description string, status string, id int, assigneeID int) error {
	if title == "" {
		return ErrTaskTitleEmpty
	}

	return s.repo.Update(ctx, title, description, status, id, assigneeID)
}
