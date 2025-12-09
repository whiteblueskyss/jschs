package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/whiteblueskyss/jschs/backend/internal/model"
)

// TeacherRepo defines DB operations for teachers.
type TeacherRepo interface {
	// Create inserts a new teacher and returns the created teacher (without exposing password_hash).
	Create(ctx context.Context, t *model.Teacher) (*model.Teacher, error)

	// GetByID fetches a teacher by ID.
	GetByID(ctx context.Context, id uuid.UUID) (*model.Teacher, error)

	// GetAll returns all teachers (with pagination added later if needed).
	GetAll(ctx context.Context) ([]*model.Teacher, error)

	// Update updates an existing teacher record and returns the updated teacher.
	Update(ctx context.Context, t *model.Teacher) (*model.Teacher, error)

	// Delete removes a teacher by ID.
	Delete(ctx context.Context, id uuid.UUID) error
}
