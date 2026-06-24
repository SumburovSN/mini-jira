package service

import (
	"context"
	"errors"
	"mini-jira/project-service/internal/model"
	"mini-jira/project-service/internal/repository"
)

type MemberService struct {
	repo *repository.MemberRepository
}

func NewMemberService(r *repository.MemberRepository) *MemberService {
	return &MemberService{repo: r}
}

func (s *MemberService) Create(ctx context.Context, projectID int, userID int, role model.Role) (*model.Member, error) {
	m := &model.Member{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}

	if !role.IsValid() {
		return nil, ErrInvalidRole
	}

	err := s.repo.Create(ctx, m)
	return m, err
}

func (s *MemberService) List(ctx context.Context, projectID int) ([]model.Member, error) {
	return s.repo.GetAllByProject(ctx, projectID)
}

func (s *MemberService) Delete(ctx context.Context, memberID int) error {
	return s.repo.Delete(ctx, memberID)
}

func (s *MemberService) GetById(ctx context.Context, id int) (*model.Member, error) {
	m, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrMemberNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}
	return m, nil
}

func (s *MemberService) UpdateRole(ctx context.Context, memberID int, role model.Role) error {
	if role == "" {
		return ErrRoleEmpty
	}

	if !role.IsValid() {
		return ErrInvalidRole
	}

	return s.repo.UpdateRole(ctx, memberID, role)
}

func (s *MemberService) GetUserRole(ctx context.Context, projectID int, userID int) (model.Role, error) {
	return s.repo.GetUserRole(ctx, projectID, userID)
}

var ErrMemberNotFound = errors.New("member not found")
var ErrRoleEmpty = errors.New("role is empty")
var ErrInvalidRole = errors.New("invalid role")
