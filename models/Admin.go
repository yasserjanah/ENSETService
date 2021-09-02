package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Admin is used by pop to map your admins database table to your go code.
type Admin struct {
	ID                   uuid.UUID    `json:"id" db:"id"`
	Email                string       `json:"email" db:"email"`
	PasswordHash         string       `json:"-" db:"password_hash"`
	Password             nulls.String `json:"-" db:"-"`
	PasswordConfirmation nulls.String `json:"-" db:"-"`
	CreatedAt            time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at" db:"updated_at"`
}

// Create new Admin (TODO: find a way to remove registration)
func (a *Admin) Create(tx *pop.Connection) (*validate.Errors, error) {
	a.Email = strings.ToLower(a.Email)
	pb, err := a.Password.MarshalJSON()
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	ph, err := bcrypt.GenerateFromPassword([]byte(pb), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	a.PasswordHash = string(ph)
	return tx.ValidateAndCreate(a)
}

// String is not required by pop and may be deleted
func (a Admin) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Admins is not required by pop and may be deleted
type Admins []Admin

// String is not required by pop and may be deleted
func (a Admins) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Admin) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Email, Name: "Email"},
		&validators.StringIsPresent{Field: a.PasswordHash, Name: "PasswordHash"},
		&validators.StringIsPresent{Field: a.Password.String, Name: "Password"},
		&validators.StringIsPresent{Field: a.PasswordConfirmation.String, Name: "PasswordConfirmation"},
		&validators.FuncValidator{
			Field:   a.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", a.Email)
				if a.ID != uuid.Nil {
					q = q.Where("id != ?", a.ID)
				}
				b, err := q.Exists(a)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Admin) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Password.String, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: a.Password.String, Field2: a.PasswordConfirmation.String, Message: "Password does not match confirmation"},
	), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Admin) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
