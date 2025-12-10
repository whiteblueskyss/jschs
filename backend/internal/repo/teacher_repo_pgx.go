package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/whiteblueskyss/jschs/backend/internal/model"
)

type teacherRepo struct {
	db *pgxpool.Pool
}

// Constructor: used by service layer to create a repo
func NewTeacherRepo(db *pgxpool.Pool) TeacherRepo {
	return &teacherRepo{db: db}
}

// --- Methods will be implemented step-by-step ---

func (r *teacherRepo) Create(ctx context.Context, t *model.Teacher) (*model.Teacher, error) {
	query := `
		INSERT INTO teachers (
			email, password_hash, full_name, phone,
			is_active, photo, date_of_birth, joining_date,
			gender, bio, address, designation, qualification
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, NULLIF($7, '')::date, NULLIF($8, '')::date,
			$9, $10, $11, $12, $13 
		)
		RETURNING id;
	`

	var id uuid.UUID
	err := r.db.QueryRow(
		ctx,
		query,
		t.Email,
		t.PasswordHash,
		t.FullName,
		t.Phone,
		t.IsActive,
		t.Photo,
		t.DateOfBirth,
		t.JoiningDate,
		t.Gender,
		t.Bio,
		t.Address,
		t.Designation,
		t.Qualification,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	t.ID = id
	return t, nil
}

func (r *teacherRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Teacher, error) {
	return nil, nil
}

func (r *teacherRepo) GetAll(ctx context.Context) ([]*model.Teacher, error) {
	return nil, nil
}

func (r *teacherRepo) Update(ctx context.Context, t *model.Teacher) (*model.Teacher, error) {
	return nil, nil
}

func (r *teacherRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
