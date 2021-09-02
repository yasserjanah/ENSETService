package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Student is used by pop to map your students database table to your go code.
type Student struct {
	ID          uuid.UUID `json:"id" db:"id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	AvatarURL   string    `json:"avatar_url" db:"avatar_url"`
	Ip          string    `json:"ip" db:"ip"`
	UserAgent   string    `json:"user_agent" db:"user_agent"`
	Provider    string    `json:"provider" db:"provider"`
	ProviderID  string    `json:"provider_id" db:"provider_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Create new Student
func (s *Student) Create(tx *pop.Connection) error {
	return tx.Create(s)
}

// String is not required by pop and may be deleted
func (s Student) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Students is not required by pop and may be deleted
type Students []Student

// String is not required by pop and may be deleted
func (s Students) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Student) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		/*		&validators.StringIsPresent{Field: s.FirstName, Name: "FirstName"},
				&validators.StringIsPresent{Field: s.LastName, Name: "LastName"},*/
		&validators.StringIsPresent{Field: s.Email, Name: "Email"},
		/*		&validators.StringIsPresent{Field: s.PhoneNumber, Name: "PhoneNumber"},
		 */&validators.StringIsPresent{Field: s.AvatarURL, Name: "AvatarURL"},
		&validators.StringIsPresent{Field: s.Ip, Name: "Ip"},
		&validators.StringIsPresent{Field: s.UserAgent, Name: "UserAgent"},
		&validators.StringIsPresent{Field: s.Provider, Name: "Provider"},
		&validators.StringIsPresent{Field: s.ProviderID, Name: "ProviderID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Student) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Student) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
