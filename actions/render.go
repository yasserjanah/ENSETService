package actions

import (
	"ensetservice/models"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.plush.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			"ParseTime": func(t time.Time) string {
				return t.Format("Mon Jan _2 2006 15:04:05")
			},
			"GetStudent": func(sid string) *models.Student {
				s := &models.Student{}
				tx := models.DB
				err := tx.Where("id = ?", sid).First(s)
				if err != nil {
					return nil
				}
				return s
			},
			"UserFromEmail": func(email string) string {
				return strings.Split(email, "@")[0]
			},
			"SplitDocPath": func(docPath string) string {
				return strings.Split(docPath, "__")[2]
			},
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
		},
	})
}
