package actions

import (
	"ensetservice/models"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {

	if uid := c.Session().Get("enset_student_id"); uid != nil {
		docs := &models.Documents{}
		tx := c.Value("tx").(*pop.Connection)
		tx.Where("student_id = ?", uid.(uuid.UUID).String()).Order("created_at desc").All(docs)

		fmt.Println()
		fmt.Println(docs)
		fmt.Println()
		c.Set("docs", docs)
	}

	if uid := c.Session().Get("enset_admin_id"); uid != nil {
		docs := &models.Documents{}
		tx := c.Value("tx").(*pop.Connection)
		tx.Order("created_at desc").All(docs)

		fmt.Println()
		fmt.Println(docs)
		fmt.Println()
		c.Set("docs", docs)
	}

	return c.Render(http.StatusOK, r.HTML("index.html"))
}
