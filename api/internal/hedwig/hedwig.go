package hedwig

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	ttmpl "text/template"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/util"
)

const (
	// List of templates. All of these must be parsed, validated in NewHedwig()
	InviteEmployee = "invite-employee"
)

type Hedwig interface {
	GenerateEmail(req GenerateEmailReq) (db.Email, error)
}

type hedwig struct {
	util.Logger
}

func NewHedwig(log util.Logger) (Hedwig, error) {
	for _, tmpl := range []string{InviteEmployee} {
		fi, err := os.Stat(filepath.Join("hedwig", "templates", tmpl+".txt"))
		if err != nil {
			return nil, err
		}
		if fi.Size() == 0 {
			return nil, fmt.Errorf("template %s is empty", tmpl)
		}
	}
	return &hedwig{Logger: log}, nil
}

type GenerateEmailReq struct {
	TemplateName string
	Args         map[string]string

	EmailFrom string
	EmailTo   []string
	EmailCC   []string
	Subject   string
}

func (h *hedwig) GenerateEmail(req GenerateEmailReq) (db.Email, error) {
	plainText, err := os.ReadFile(
		filepath.Join("hedwig", "templates", req.TemplateName+".txt"),
	)
	if err != nil {
		return db.Email{}, err
	}

	var textBody bytes.Buffer
	err = ttmpl.Must(
		ttmpl.New("text").Parse(string(plainText)),
	).Execute(&textBody, req.Args)
	if err != nil {
		h.Err("text template", "template", req.TemplateName, "error", err)
		return db.Email{}, err
	}

	htmlText, err := os.ReadFile(
		filepath.Join("hedwig", "templates", req.TemplateName+".html"),
	)
	if err != nil {
		h.Err("html template", "template", req.TemplateName, "error", err)
		return db.Email{}, err
	}

	var htmlBody bytes.Buffer
	err = ttmpl.Must(
		ttmpl.New("html").Parse(string(htmlText)),
	).Execute(&htmlBody, req.Args)
	if err != nil {
		h.Err("html template", "template", req.TemplateName, "error", err)
		return db.Email{}, err
	}

	return db.Email{
		EmailFrom:     req.EmailFrom,
		EmailTo:       req.EmailTo,
		EmailCC:       req.EmailCC,
		EmailSubject:  req.Subject,
		EmailTextBody: textBody.String(),
		EmailHTMLBody: htmlBody.String(),
	}, nil
}
