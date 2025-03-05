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
	InviteEmployee               = "invite-employee"
	InviteHubUser                = "invite-hub-user"
	HubUserTFA                   = "hub-user-tfa"
	HubPasswordReset             = "hub-password-reset"
	ShortlistApplication         = "shortlist-application"
	RejectApplication            = "reject-application"
	NotifyNewInterviewer         = "notify-new-interviewer"
	NotifyWatchersNewInterviewer = "notify-watchers-new-interviewer"
	RemovedInterviewerNotify     = "removed-interviewer-notify"
	NotifyApplicantInterview     = "notify-applicant-interview"
	NotifyCandidateOffer         = "notify-candidate-offer"
	AddOfficialEmail             = "add-official-email"
	EndorsementRequest           = "endorsement-request"
)

type Hedwig interface {
	GenerateEmail(req GenerateEmailReq) (db.Email, error)
}

type hedwig struct {
	util.Logger
}

func NewHedwig(log util.Logger) (Hedwig, error) {
	for _, tmpl := range []string{
		InviteEmployee,
		HubUserTFA,
		HubPasswordReset,
		ShortlistApplication,
		RejectApplication,
		NotifyNewInterviewer,
		NotifyWatchersNewInterviewer,
		RemovedInterviewerNotify,
		NotifyApplicantInterview,
		NotifyCandidateOffer,
		AddOfficialEmail,
		EndorsementRequest,
	} {
		fi, err := os.Stat(filepath.Join("hedwig", "templates", tmpl+".txt"))
		if err != nil {
			return nil, err
		}
		if fi.Size() == 0 {
			return nil, fmt.Errorf("template %s.txt is empty", tmpl)
		}

		fi, err = os.Stat(filepath.Join("hedwig", "templates", tmpl+".html"))
		if err != nil {
			return nil, err
		}
		if fi.Size() == 0 {
			return nil, fmt.Errorf("template %s.html is empty", tmpl)
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
	// TODO: We need to get a preferred_lang parameter and use the appropriate template. Right now our email templates are hardcoded on English. We need to support other languages.
	plainText, err := os.ReadFile(
		filepath.Join("hedwig", "templates", req.TemplateName+".txt"),
	)
	if err != nil {
		h.Err("Failed to read template", "error", err)
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
		EmailState:    db.EmailStatePending,
	}, nil
}
