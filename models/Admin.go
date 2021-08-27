package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"time"
	"github.com/gobuffalo/validate/v3/validators"
)
// Admin is used by pop to map your admins database table to your go code.
type Admin struct {
    ID uuid.UUID `json:"id" db:"id"`
    Email string `json:"email" db:"email"`
    PasswordHash string `json:"password_hash" db:"password_hash"`
    Password string `json:"password" db:"password"`
    PasswordConfirmation string `json:"password_confirmation" db:"password_confirmation"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Create new Admin (TODO: find a way to remove registration)
func (a *Admin) Create(tx *pop.Connection) (*validate.Errors, error) {
	a.Email = strings.ToLower(a.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
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
		&validators.StringIsPresent{Field: a.Password, Name: "Password"},
		&validators.StringIsPresent{Field: a.PasswordConfirmation, Name: "PasswordConfirmation"},
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				b, err = q.Exists(u)
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
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Admin) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
