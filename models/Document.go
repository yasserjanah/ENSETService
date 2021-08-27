package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"time"
	"github.com/gobuffalo/validate/v3/validators"
)
// Document is used by pop to map your documents database table to your go code.
type Document struct {
    ID uuid.UUID `json:"id" db:"id"`
    DocName string `json:"doc_name" db:"doc_name"`
    StudentID string `json:"student_id" db:"student_id"`
    IsDone bool `json:"is_done" db:"is_done"`
    Status string `json:"status" db:"status"`
    DocPath string `json:"doc_path" db:"doc_path"`
    Message string `json:"message" db:"message"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Create new Document
func (d *Document) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(d)
}

// String is not required by pop and may be deleted
func (d Document) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Documents is not required by pop and may be deleted
type Documents []Document

// String is not required by pop and may be deleted
func (d Documents) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *Document) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.DocName, Name: "DocName"},
		&validators.StringIsPresent{Field: d.StudentID, Name: "StudentID"},
		&validators.StringIsPresent{Field: d.Status, Name: "Status"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *Document) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *Document) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
