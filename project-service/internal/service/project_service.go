package service

import (
	"context"
	"errors"
	"mini-jira/project-service/internal/model"
	"mini-jira/project-service/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(r *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: r}
}

func (s *ProjectService) Create(ctx context.Context, name string, ownerID int) (*model.Project, error) {
	p := &model.Project{
		Name:    name,
		OwnerID: ownerID,
	}
	err := s.repo.Create(ctx, p)
	return p, err
}

func (s *ProjectService) List(ctx context.Context, ownerID int) ([]model.Project, error) {
	return s.repo.GetAllByOwner(ctx, ownerID)
}

func (s *ProjectService) Delete(ctx context.Context, id int, ownerID int) error {
	return s.repo.Delete(ctx, id, ownerID)
}

var ErrProjectNotFound = errors.New("project not found")

func (s *ProjectService) GetById(ctx context.Context, id int, ownerID int) (*model.Project, error) {
	p, err := s.repo.GetById(ctx, id, ownerID)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return p, nil
}
