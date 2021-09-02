package actions

import (
	"ensetservice/models"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

// Middlewares

// SetCurrentUser attempts to find a user based on the enset_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {

		if auid := c.Session().Get("enset_admin_id"); auid != nil {
			a := &models.Admin{}
			tx := c.Value("tx").(*pop.Connection)
			q := tx.Where("id = ?", auid)
			ae, err := q.Exists(a)
			if ae {
				// fmt.Println("*** Setting enset_admin...")
				err := q.First(a)
				if err != nil {
					return errors.WithStack(err)
				}
				c.Set("enset_admin", a)
			}
			if err != nil {
				return errors.WithStack(err)
			}
			fmt.Println(a)

		}

		if suid := c.Session().Get("enset_student_id"); suid != nil {
			s := &models.Student{}
			tx := c.Value("tx").(*pop.Connection)
			q := tx.Where("id = ?", suid)
			se, err := q.Exists(s)
			if se {
				// fmt.Println("*** Setting enset_student...")
				err := q.First(s)
				if err != nil {
					return errors.WithStack(err)
				}
				c.Set("enset_student", s)
			}
			if err != nil {
				return errors.WithStack(err)
			}
			fmt.Println(s)

		}

		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {

		uuid := c.Session().Get("enset_admin_id")
		suid := c.Session().Get("enset_student_id")

		if uuid == nil && suid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}

		return next(c)
	}
}

// Make sure that only students have access
func StudentsOnly(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {

		if suid := c.Session().Get("enset_student_id"); suid == nil {
			c.Set("error_message", "You don't have the right to access this page (you are not Student)")
			return c.Render(http.StatusUnauthorized, r.HTML("errors/errors.html"))
		}

		return next(c)
	}
}

// Make sure that only Admins have access
func AdminsOnly(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {

		if auid := c.Session().Get("enset_admin_id"); auid == nil {
			c.Set("error_message", "You don't have the right to access this page (you are not Admin)")
			return c.Render(http.StatusUnauthorized, r.HTML("errors/errors.html"))
		}

		return next(c)
	}
}

// Only logged OUT users can access
func Unauthorized(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		auid := c.Session().Get("enset_admin_id")
		suid := c.Session().Get("enset_student_id")
		if auid != nil || suid != nil {
			c.Set("error_message", "You're already Logged In")
			return c.Render(http.StatusUnauthorized, r.HTML("errors/errors.html"))
		}
		return next(c)
	}
}
