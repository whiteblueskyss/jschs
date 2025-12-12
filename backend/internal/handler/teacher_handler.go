package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/whiteblueskyss/jschs/backend/internal/model"
	"github.com/whiteblueskyss/jschs/backend/internal/service"
)

// teacherHandler holds service dependency
type teacherHandler struct {
	svc service.TeacherService
}

// NewTeacherHandler constructs a teacher handler
func NewTeacherHandler(s service.TeacherService) *teacherHandler {
	return &teacherHandler{svc: s}
}

var validate = validator.New()

// create request payload
// i will use password "password"
type createTeacherRequest struct {
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=4"`
	FullName      string `json:"full_name" validate:"required,min=2"`
	Phone         string `json:"phone" validate:"required"`
	IsActive      *bool  `json:"is_active,omitempty"`
	Photo         string `json:"photo,omitempty"`
	DateOfBirth   string `json:"date_of_birth,omitempty" validate:"omitempty,datetime=2006-01-02"`
	JoiningDate   string `json:"joining_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Gender        string `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Bio           string `json:"bio,omitempty"`
	Address       string `json:"address,omitempty"`
	Designation   string `json:"designation,omitempty"`
	Qualification string `json:"qualification,omitempty"`
}

// RegisterHandler handles POST /api/v1/teachers
func (h *teacherHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req createTeacherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}
	// validation

	if err := validate.Struct(req); err != nil {
		// build a sensible message
		if ve, ok := err.(validator.ValidationErrors); ok {
			// return first validation error message
			http.Error(w, ve.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "validation error", http.StatusBadRequest)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// build model
	teacher := &model.Teacher{
		// ID will be set by DB
		Email:         req.Email,
		FullName:      req.FullName,
		Phone:         req.Phone,
		IsActive:      isActive,
		Photo:         req.Photo,
		DateOfBirth:   req.DateOfBirth,
		JoiningDate:   req.JoiningDate,
		Gender:        req.Gender,
		Bio:           req.Bio,
		Address:       req.Address,
		Designation:   req.Designation,
		Qualification: req.Qualification,
	}

	created, err := h.svc.Register(r.Context(), teacher, req.Password)
	if err != nil {
		http.Error(w, "failed to create teacher: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// created.PasswordHash is cleared by service; still remove defensively
	created.PasswordHash = ""
	_ = json.NewEncoder(w).Encode(created)
}

// GET /teachers/{id}
func (h *teacherHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	t, err := h.svc.Get(r.Context(), id)
	if err != nil {
		// if not found, return 404 for common pgx ErrNoRows
		// compare error string or better: import pgx.ErrNoRows if you want exact check
		http.Error(w, "teacher not found", http.StatusNotFound)
		return
	}
	if t == nil {
		http.Error(w, "teacher not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t)
}

// GET /teachers
func (h *teacherHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch teachers: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

// PUT /teachers/{id}
func (h *teacherHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req model.Teacher
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}

	// ensure the ID in path is authoritative
	req.ID = id

	// Do not allow password change via this route.
	req.PasswordHash = "" // ensure repo doesn't replace password unless service calls ChangePassword.

	updated, err := h.svc.UpdateProfile(r.Context(), &req)
	if err != nil {
		http.Error(w, "failed to update teacher: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updated)
}

// small helper to register routes in router
func (h *teacherHandler) Routes(r chi.Router) {
	r.Post("/teachers", h.RegisterHandler)
	r.Get("/teachers/{id}", h.GetByID)
	r.Get("/teachers", h.GetAll)
	r.Put("/teachers/{id}", h.Update)
}

// If need GetByID later: have to use chi.URLParam(r, "id") and uuid.Parse
func getIDParam(r *http.Request, name string) (uuid.UUID, error) {
	val := chi.URLParam(r, name)
	return uuid.Parse(val)
}
