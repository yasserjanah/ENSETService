package actions

import (
	"ensetservice/mailers"
	"ensetservice/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

var AvailableDocuments = []string{"Relevé de notes", "certificat scolaire"}

// DocumentNew default implementation.
func DocumentNew(c buffalo.Context) error {
	uid := c.Session().Get("enset_student_id")
	s := &models.Student{}
	tx := c.Value("tx").(*pop.Connection)
	tx.Where("id = ?", uid.(uuid.UUID).String()).First(s)
	c.Set("student", s)
	c.Set("av_docs", AvailableDocuments)
	return c.Render(http.StatusOK, r.HTML("documents/new.html"))
}

// DocumentCreate add new document.
func DocumentCreate(c buffalo.Context) error {
	// make sure that only students has the permission to add new document

	d := &models.Document{}
	if err := c.Bind(d); err != nil {
		return errors.WithStack(err)
	}

	found := false
	for _, c := range AvailableDocuments {
		if c == d.DocName {
			found = true
		}
	}

	if !found {
		c.Set("errors", errors.New("Invalid document Name !! "))
		return c.Render(http.StatusBadRequest, r.HTML("documents/new.html"))
	}

	// set the student_id (don't let the user set it, TO AVOID IDOR)
	uid := c.Session().Get("enset_student_id")
	d.StudentID = uid.(uuid.UUID).String()

	// set document status to "p" (pending)
	d.Status = "PENDING"

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := d.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	// handle student personel information
	s := &models.Student{}
	tx.Where("id = ?", uid.(uuid.UUID).String()).First(s)

	if err := c.Bind(s); err != nil {
		return errors.WithStack(err)
	}
	tx.Save(s)

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusBadRequest, r.HTML("documents/new.html"))
	}

	// fmt.Println()
	// fmt.Println("*** adding new document:", d)
	// fmt.Println()

	c.Flash().Add("success", "Votre demande a été créé avec succès.!")

	// send email to the student
	go func() {

		err = mailers.SendNotification([]string{s.Email}, "Demande créée", `Votre demande a été créée avec succès.
		nous vous tiendrons au courant de l’état de votre demande.`)

		if err != nil {
			fmt.Println("*** error sending email:", err)
		}

	}()

	// send email to the admin
	go func() {
		admins := []*models.Admin{}
		tx.All(admins)

		for _, a := range admins {
			err = mailers.SendNotification([]string{a.Email}, "Nouvelle Demande", fmt.Sprintf(`Une nouvelle demande a été créé.
			par l'étudiant %v %v`, s.FirstName, s.LastName))

			if err != nil {
				fmt.Println("*** error sending email:", err)
			}
		}

	}()

	return c.Redirect(http.StatusFound, "/")
}

// Process documents
func DocumentProcessorNew(c buffalo.Context) error {
	// Get document ID from URL
	docID := c.Param("docID")

	tx := c.Value("tx").(*pop.Connection)

	// find document with the ID
	doc := &models.Document{}
	q := tx.Where("id = ?", docID)

	ed, err := q.Exists(doc)
	if err != nil {
		return errors.WithStack(err)
	}

	// if document with that ID not exists
	if !ed {
		c.Flash().Add("danger", "can't find any Document")
		c.Set("hide_form", true)
		return c.Render(http.StatusOK, r.HTML("documents/process.html"))
	}

	err = q.First(doc)
	if err != nil {
		return errors.WithStack(err)
	}
	// do the check if the document is already resolved
	if doc.IsDone {
		if doc.Status == "REJECTED" {
			c.Flash().Add("danger", "Document Already Rejected !!")
		} else {
			c.Flash().Add("warning", "Document Already Processed !!")
		}

	}

	c.Set("doc", doc)
	return c.Render(http.StatusOK, r.HTML("documents/process.html"))
}

func DocumentProcessorCreate(c buffalo.Context) error {
	// Get document ID from URL
	docID := c.Param("docID")
	c.Set("docID", docID)

	tx := c.Value("tx").(*pop.Connection)

	doc := &models.Document{}
	q := tx.Where("id = ?", docID)

	ed, err := q.Exists(doc)
	if err != nil {
		return errors.WithStack(err)
	}

	// if document with that ID not exists
	if !ed {
		c.Flash().Add("danger", "can't find any Document")
		c.Set("hide_form", true)
		return c.Render(http.StatusOK, r.HTML("documents/process.html"))
	}

	// find document with the ID
	err = q.First(doc)
	if err != nil {
		return errors.WithStack(err)
	}

	var Status map[string]string = map[string]string{"Reject": "REJECTED", "Update": "DONE", "Upload": "DONE"}
	doc.Status = Status[strings.TrimSpace(c.Request().FormValue("status"))]

	if doc.Status != "REJECTED" {

		// var Status map[string]string = map[string]string{"Reject": "REJECTED", "Update": "DONE", "Upload": "DONE"}
		// doc.Status = Status[strings.TrimSpace(c.Request().FormValue("status"))]

		f, err := c.File("pdoc")
		if err == nil {
			docPath := fmt.Sprintf("doc__%s__%s", docID, f.Filename)
			if f.Valid() {
				err := saveDocument(docPath, f)
				if err != nil {
					return errors.WithStack(err)
				}
			}

			doc.DocPath = docPath
			doc.IsDone = true
		} else {
			doc.Status = "PENDING"
			if doc.IsDone {
				doc.Status = "DONE"
			}

		}
	}

	if doc.Status == "REJECTED" {
		doc.IsDone = true
	}

	if err := c.Bind(doc); err != nil {
		fmt.Println("BIND: ", err)
		return errors.WithStack(err)
	}

	doc.Message = strings.TrimSpace(doc.Message)

	err = tx.Save(doc)

	// fmt.Println()
	// fmt.Println(doc)
	// fmt.Println()

	if err != nil {
		fmt.Println("SAVE: ", err)
		return errors.WithStack(err)
	}

	c.Set("doc", doc)

	uid := doc.StudentID
	s := &models.Student{}
	tx.Where("id = ?", uid).First(s)

	if doc.Status != "REJECTED" {
		// do the check if the document is already resolved
		if doc.IsDone {
			c.Flash().Add("warning", "Document Already processed !!")
		}

		// send email to the student
		go func() {

			err = mailers.SendNotification([]string{s.Email}, "Document est prêt", `Votre demande est bien traité, vous pouvez télécharger votre document dès maintenant.`)

			if err != nil {
				fmt.Println("*** error sending email:", err)
			}

		}()
		c.Flash().Add("success", "Document Updated succesfully !!")
		return c.Render(http.StatusOK, r.HTML("documents/process.html"))
	}

	// send email to the student
	go func() {

		err = mailers.SendNotification([]string{s.Email}, "Demande rejetée", `Votre demande est rejetée, pour plus de détails veuillez consulter votre compte.`)

		if err != nil {
			fmt.Println("*** error sending email:", err)
		}

	}()

	c.Flash().Add("success", "Document Rejected succesfully !!")
	return c.Redirect(http.StatusFound, "/")
}

func saveDocument(path string, f binding.File) error {
	dir := filepath.Join(".", "uploads/documents/")

	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithStack(err)
	}

	cf, err := os.Create(filepath.Join(dir, path))
	if err != nil {
		return errors.WithStack(err)
	}

	defer cf.Close()
	_, err = io.Copy(cf, f)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Student can Download document
func DocumentDownloader(c buffalo.Context) error {
	// Get document ID from URL
	docID := c.Param("docID")

	tx := c.Value("tx").(*pop.Connection)

	doc := &models.Document{}

	// TODO : make sure student can only downloads documents he own
	uid := c.Session().Get("enset_student_id")
	q := tx.Where("id = ? AND student_id = ?", docID, uid.(uuid.UUID).String())

	ed, err := q.Exists(doc)
	if err != nil {
		return errors.WithStack(err)
	}

	// if document with that ID not exists
	if !ed {
		c.Flash().Add("danger", "can't find any Document")
		return c.Render(http.StatusOK, r.HTML("index.html"))
	}

	// find document with the ID
	err = q.First(doc)
	if err != nil {
		return errors.WithStack(err)
	}

	// do the check if the document is already resolved
	// TODO: check that the doc is not rejected
	if doc.IsDone && doc.Status == "DONE" {
		fmt.Println()
		fmt.Println("Document is Downloadable !!")
		fmt.Println()
		path := filepath.Join(".", "uploads/documents/")
		docPath := filepath.Join(path, doc.DocPath)
		ds, err := os.Open(docPath)
		if err != nil {
			return errors.WithStack(err)
		}
		return c.Render(http.StatusOK, r.Download(c, strings.Split(doc.DocPath, "__")[2], ds))

	}

	return c.Render(http.StatusBadRequest, r.String(fmt.Sprintf("Document status is : %s - Nice try (^_^)", doc.Status)))
}
