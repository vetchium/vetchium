package granger

import (
	"bytes"
	"context"
	"html/template"
	ttmpl "text/template"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/pkg/vetchi"
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

func (g *Granger) createOnboardEmails(quit chan struct{}) {
	defer g.wg.Done()

	for {
		select {
		case <-quit:
			g.log.Debug("createOnboardEmails quitting")
			return
		case <-time.After(5 * time.Second):
			ctx := context.Background()
			onboardInfo, err := g.db.DeQOnboard(ctx)
			if err != nil {
				continue
			}

			if onboardInfo == nil {
				g.log.Debug("no pending employer onboard email generation")
				continue
			}

			g.log.Info("onboard invites", "onboardInfo", onboardInfo)

			// TODO: Should we read the length from a config?
			token := util.RandomString(vetchi.InviteTokenLenBytes)

			link := vetchi.EmployerBaseURL + "/onboard/" + token

			var textBody bytes.Buffer
			err = ttmpl.Must(
				ttmpl.New("text").Parse(textMailTemplate),
			).Execute(&textBody, map[string]string{
				"Domain": onboardInfo.DomainName,
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
				"Domain": onboardInfo.DomainName,
				"Link":   link,
			})
			if err != nil {
				g.log.Error("email html template failed", "error", err)
				continue
			}

			email := db.Email{
				EmailFrom:     vetchi.EmailFrom,
				EmailTo:       []string{onboardInfo.AdminEmailAddr},
				EmailSubject:  subject,
				EmailHTMLBody: htmlBody.String(),
				EmailTextBody: textBody.String(),
			}

			// Errors are already logged, so we can ignore the return value
			_ = g.db.CreateOnboardEmail(
				ctx,
				db.OnboardEmailInfo{
					EmployerID:         onboardInfo.EmployerID,
					OnboardSecretToken: token,
					TokenValidMins:     g.onboardTokenLife.Minutes(),
					Email:              email,
				},
			)
		}
	}
}
