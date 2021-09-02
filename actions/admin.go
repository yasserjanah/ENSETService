package actions

import (
	"database/sql"
	"ensetservice/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AdminNew default implementation.
func AdminNew(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("admin/login.plush.html"))
}

// AdminLogin default implementation.
func AdminLogin(c buffalo.Context) error {
	// TODO: check if the user already has a session

	a := &models.Admin{}
	if err := c.Bind(a); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user with the email
	err := tx.Where("email = ?", strings.ToLower(strings.TrimSpace(a.Email))).First(a)

	// helper function to handle bad attempts
	bad := func() error {
		verrs := validate.NewErrors()
		verrs.Add("email", "Invalid Email / Password !!")
		c.Flash().Add("danger", "Invalid Email / Password !!")

		c.Set("errors", verrs)

		return c.Render(http.StatusUnauthorized, r.HTML("admin/login.plush.html"))
	}

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with the supplied email address.
			return bad()
		}
		return errors.WithStack(err)
	}

	// confirm that the given password matches the hashed password from the db
	pb, err := a.Password.MarshalJSON()
	err = bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(pb))
	if err != nil {
		return bad()
	}
	fmt.Println(a)
	c.Session().Set("enset_admin_id", a.ID)
	c.Flash().Add("success", "Welcome Back to Buffalo!")

	redirectURL := "/"
	if redir, ok := c.Session().Get("redirectURL").(string); ok && redir != "" {
		redirectURL = redir
	}

	fmt.Println(redirectURL)

	return c.Redirect(http.StatusFound, redirectURL)
}

// AdminLogout default implementation.
func AdminLogout(c buffalo.Context) error {
	if auid := c.Session().Get("enset_admin_id"); auid != nil {
		c.Session().Clear()
		err := c.Session().Save()
		if err != nil {
			return errors.WithStack(err)
		}
		c.Flash().Add("success", "You have been logged out!")
	}
	// TODO: make sure logged user only has access
	return c.Redirect(http.StatusFound, "/")
}

// TODO: to be removed (don't allow admin registration)
// // AdminRegisterNew renders the admin registration form
// func AdminRegisterNew(c buffalo.Context) error {
// 	a := models.Admin{}
// 	c.Set("admin", a)
// 	return c.Render(200, r.HTML("admin/new.plush.html"))
// }

// // AdminRegisterCreate registers a new admin with the application.
// func AdminRegisterCreate(c buffalo.Context) error {
// 	a := &models.Admin{}
// 	if err := c.Bind(a); err != nil {
// 		return errors.WithStack(err)
// 	}

// 	tx := c.Value("tx").(*pop.Connection)
// 	verrs, err := a.Create(tx)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	if verrs.HasAny() {
// 		c.Set("admin", a)
// 		c.Set("errors", verrs)
// 		return c.Render(200, r.HTML("admin/new.plush.html"))
// 	}

// 	c.Session().Set("enset_admin_id", a.ID)
// 	c.Flash().Add("success", "Welcome to ENSET Service!")

// 	return c.Redirect(302, "/")
// }
