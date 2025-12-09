package model

import "github.com/google/uuid"

// Teacher represents the teacher entity mapped to the teachers table.
type Teacher struct {
	ID            uuid.UUID `json:"id"` // primary key
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"` // omit in JSON output
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	IsActive      bool      `json:"is_active"`
	Photo         string    `json:"photo"`
	DateOfBirth   string    `json:"date_of_birth"` // use ISO date "YYYY-MM-DD" for JSON; DB will use date type
	JoiningDate   string    `json:"joining_date"`
	Gender        string    `json:"gender"`
	Bio           string    `json:"bio"`
	Address       string    `json:"address"`
	Designation   string    `json:"designation"`
	Qualification string    `json:"qualification"`
}
