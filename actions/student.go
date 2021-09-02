package actions

import (
	"ensetservice/models"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/pkg/errors"
)

func init() {
	gothic.Store = App().SessionStore
	goth.UseProviders(google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/students/login/google/callback")))
}

// StudentLogin default implementation.
func StudentLogin(c buffalo.Context) error {
	// TODO: check if the user already has a session

	gu, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}
	tx := c.Value("tx").(*pop.Connection)

	//check if it is an enset-media.ac.ma email
	if domain := strings.Split(gu.Email, "@"); domain[1] != "enset-media.ac.ma" {
		c.Set("error_message", fmt.Sprintf("%v is not an ENSETM email", gu.Email))
		return c.Render(http.StatusUnauthorized, r.HTML("errors/errors.html"))
	}

	//check if user already exist
	q := tx.Where("provider = ? AND provider_id = ? AND email = ?", gu.Provider, gu.UserID, gu.Email)
	exists, err := q.Exists("students")
	if err != nil {
		return errors.WithStack(err)
	}

	// Creating student
	var student *models.Student = new(models.Student)
	if exists {
		err = q.First(student)
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Println()
		fmt.Println("Found existing user with this email")
		fmt.Println()
		// update student
		student.UserAgent = c.Request().UserAgent()
		student.Ip = c.Request().RemoteAddr
		tx.Save(student)
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		fmt.Println()
		fmt.Println("New Email has been registered")
		fmt.Println()
		student.FirstName = gu.FirstName
		student.LastName = gu.LastName
		student.Email = gu.Email
		student.AvatarURL = gu.AvatarURL
		student.UserAgent = c.Request().UserAgent()
		student.Ip = c.Request().RemoteAddr
		student.Provider = gu.Provider
		student.ProviderID = gu.UserID

		err := student.Create(tx)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	c.Session().Set("enset_student_id", student.ID)
	if err = c.Session().Save(); err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "Welcome Back!")

	return c.Redirect(http.StatusFound, "/")
	/*	return c.Render(http.StatusOK, r.HTML("student/login_callback.html"))
	 */
}

// StudentLogout default implementation.
func StudentLogout(c buffalo.Context) error {
	if suid := c.Session().Get("enset_student_id"); suid != nil {
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
