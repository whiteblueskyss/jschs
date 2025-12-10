package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/whiteblueskyss/jschs/backend/internal/model"
	"github.com/whiteblueskyss/jschs/backend/internal/repo"
)

// teacherService is the concrete implementation of TeacherService.
type teacherService struct {
	repo repo.TeacherRepo
}

// NewTeacherService constructs a TeacherService.
func NewTeacherService(r repo.TeacherRepo) TeacherService {
	return &teacherService{repo: r}
}

// Register hashes the plain password, stores the teacher, and returns the created teacher without password.
func (s *teacherService) Register(ctx context.Context, t *model.Teacher, plainPassword string) (*model.Teacher, error) {
	if t == nil {
		return nil, errors.New("teacher is nil")
	}
	if plainPassword == "" {
		return nil, errors.New("password is required")
	}

	// hash password using bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	t.PasswordHash = string(hashed)

	// delegate to repo to persist
	created, err := s.repo.Create(ctx, t)
	if err != nil {
		return nil, err
	}

	// remove password hash from returned object for safety (defense-in-depth)
	created.PasswordHash = ""
	return created, nil
}

func (s *teacherService) Get(ctx context.Context, id uuid.UUID) (*model.Teacher, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// defensive: hide password hash when returning to callers
	if t != nil {
		t.PasswordHash = ""
	}
	return t, nil
}

// GetAll is not implemented yet.
func (s *teacherService) GetAll(ctx context.Context) ([]*model.Teacher, error) {
	return nil, errors.New("not implemented")
}

// Authenticate is not implemented yet.
func (s *teacherService) Authenticate(ctx context.Context, email, plainPassword string) (*model.Teacher, error) {
	return nil, errors.New("not implemented")
}

// UpdateProfile is not implemented yet.
func (s *teacherService) UpdateProfile(ctx context.Context, t *model.Teacher) (*model.Teacher, error) {
	return nil, errors.New("not implemented")
}

// ChangePassword is not implemented yet.
func (s *teacherService) ChangePassword(ctx context.Context, id uuid.UUID, newPlainPassword string) error {
	return errors.New("not implemented")
}

// Delete is not implemented yet.
func (s *teacherService) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented")
}
