package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

// create request payload
type createTeacherRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	FullName      string `json:"full_name"`
	Phone         string `json:"phone"`
	IsActive      *bool  `json:"is_active,omitempty"`
	Photo         string `json:"photo,omitempty"`
	DateOfBirth   string `json:"date_of_birth,omitempty"`
	JoiningDate   string `json:"joining_date,omitempty"`
	Gender        string `json:"gender,omitempty"`
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
	// basic validation
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		http.Error(w, "email, password and full_name are required", http.StatusBadRequest)
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

// Optional: small helper to register routes in router
func (h *teacherHandler) Routes(r chi.Router) {
	r.Post("/teachers", h.RegisterHandler)
	r.Get("/teachers/{id}", h.GetByID)
}

// If need GetByID later: have to use chi.URLParam(r, "id") and uuid.Parse
func getIDParam(r *http.Request, name string) (uuid.UUID, error) {
	val := chi.URLParam(r, name)
	return uuid.Parse(val)
}
