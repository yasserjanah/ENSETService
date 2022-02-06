package actions

import (
	"ensetservice/models"
	"net/http"
	"net/mail"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {

	if uid := c.Session().Get("enset_student_id"); uid != nil {
		docs := &models.Documents{}
		tx := c.Value("tx").(*pop.Connection)
		tx.Where("student_id = ?", uid.(uuid.UUID).String()).Order("created_at desc").All(docs)

		// fmt.Println()
		// fmt.Println(docs)
		// fmt.Println()
		c.Set("docs", docs)
	}

	if uid := c.Session().Get("enset_admin_id"); uid != nil {
		docs := &models.Documents{}
		tx := c.Value("tx").(*pop.Connection)

		var searchKeyword string = c.Param("search")

		if searchKeyword != "" {
			c.Set("searchKeyword", searchKeyword)
			if isValidEmail(searchKeyword) {
				s := &models.Student{}
				err := tx.Where("email LIKE ?", string(searchKeyword)).First(s)
				if err != nil {
					return errors.WithStack(err)
				}
				err = tx.Where("student_id = ?", s.ID).Order("created_at desc").All(docs)
				if err != nil {
					return errors.WithStack(err)
				}
			} else {

				err := tx.Where("doc_name LIKE ?", string(searchKeyword)).Order("created_at desc").All(docs)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		} else {
			tx.Order("created_at desc").All(docs)
		}

		// fmt.Println()
		// fmt.Println(docs)
		// fmt.Println()
		c.Set("docs", docs)
	}

	return c.Render(http.StatusOK, r.HTML("index.html"))
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
