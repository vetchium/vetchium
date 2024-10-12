package granger

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"html/template"
	ttmpl "text/template"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/libvetchi"
)

const subject = "Welcome to Vetchi !"

const textMailTemplate = `Hi

You have been invited to set up password for managing {{ .Domain }} on Vetchi.

Please click the link below to set up your password:
{{ .Link }}

Thanks,
The Vetchi Team
`

const htmlMailTemplate = `Hi
<p>
You have been invited to set up password for managing {{ .Domain }} on Vetchi.
</p>
<p>
Please click the link below to set up your password:
</p>
<p>
<a href="{{ .Link }}">Set up your password</a>
</p>
<p>
In case the above link does not work, you can copy and paste the following URL in your browser:
</p>
<p>
{{ .Link }}
</p>
<p>
Thanks,
The Vetchi Team
</p>
`

func (g *Granger) createOnboardEmails() {
	defer g.wg.Done()

	for {
		select {
		case <-g.quit:
			return
		case <-time.After(3 * time.Minute):
			employers, err := g.db.GetUnmailedOnboardPendingEmployers()
			if err != nil {
				continue
			}

			for _, employer := range employers {
				g.log.Info("onboard invites", "employer", employer.ClientID)

				buff := make([]byte, 16)
				rand.Read(buff)

				token := hex.EncodeToString(buff)
				employer.OnboardSecretToken = pgtype.Text{
					String: token,
					Valid:  true,
				}

				link := libvetchi.BaseURL + "/onboard/" + token

				var textBody bytes.Buffer
				err := ttmpl.Must(
					ttmpl.New("text").Parse(textMailTemplate),
				).Execute(&textBody, map[string]string{
					"Domain": employer.ClientID,
					"Link":   link,
				})
				if err != nil {
					g.log.Error("email text template failed", "error", err)
					continue
				}

				var htmlBody bytes.Buffer
				err = template.Must(
					template.New("html").Parse(htmlMailTemplate),
				).Execute(&htmlBody, map[string]string{
					"Domain": employer.ClientID,
					"Link":   link,
				})
				if err != nil {
					g.log.Error("email html template failed", "error", err)
					continue
				}

				email := db.Email{
					EmailFrom:     libvetchi.EmailFrom,
					EmailTo:       []string{employer.OnboardAdmin.String},
					EmailSubject:  subject,
					EmailHTMLBody: htmlBody.String(),
					EmailTextBody: textBody.String(),
				}

				// Errors are already logged, so we can ignore the return value
				_ = g.db.CreateOnboardEmail(employer, email)
			}
		}
	}
}
