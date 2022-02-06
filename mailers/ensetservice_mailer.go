package mailers

import (
	"fmt"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
)

func SendNotification(to []string, subject string, msg string) error {
	m := mail.NewMessage()

	// fill in with your stuff:
	m.Subject = subject
	m.From = "ENSET Service <enset-service@enset-media.ac.ma>"
	m.To = to
	err := m.AddBody(r.HTML("ensetservice_mailer.html"), render.Data{"Message": msg})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return smtp.Send(m)
}
