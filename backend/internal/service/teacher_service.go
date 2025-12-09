package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/whiteblueskyss/jschs/backend/internal/model"
)

// TeacherService defines business-level operations.
type TeacherService interface {
	// Register creates a new teacher (handles hashing password, validation).
	// Returns the created teacher (with PasswordHash kept internal and not serialized).
	Register(ctx context.Context, t *model.Teacher, plainPassword string) (*model.Teacher, error)

	// Get retrieves teacher details by ID.
	Get(ctx context.Context, id uuid.UUID) (*model.Teacher, error)

	// GetAll retrieves all teachers.
	GetAll(ctx context.Context) ([]*model.Teacher, error)

	// Authenticate checks email+password and returns teacher if valid.
	Authenticate(ctx context.Context, email, plainPassword string) (*model.Teacher, error)

	// UpdateProfile updates teacher profile fields (not password) and returns updated teacher.
	UpdateProfile(ctx context.Context, t *model.Teacher) (*model.Teacher, error)

	// ChangePassword updates password with hashing.
	ChangePassword(ctx context.Context, id uuid.UUID, newPlainPassword string) error

	// Delete removes a teacher.
	Delete(ctx context.Context, id uuid.UUID) error
}
