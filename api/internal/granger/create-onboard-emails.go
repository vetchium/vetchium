package granger

import (
	"bytes"
	"context"
	"html/template"
	ttmpl "text/template"
	"time"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
)

const subject = "Welcome to Vetchium !"

const textMailTemplate = `Hi

You have been invited to set up password for managing {{ .Domain }} on Vetchi.

Please click the link below to set up your password:
{{ .Link }}

Thanks,
The Vetchium Team
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
The Vetchium Team
</p>
`

func (g *Granger) createOnboardEmails(quit chan struct{}) {
	g.log.Dbg("Starting createOnboardEmails job")
	defer g.log.Dbg("createOnboardEmails job finished")
	defer g.wg.Done()

	for {
		ticker := time.NewTicker(vetchi.CreateOnboardEmailsInterval)
		select {
		case <-quit:
			ticker.Stop()
			g.log.Dbg("createOnboardEmails quitting")
			return
		case <-ticker.C:
			ticker.Stop()
			ctx := context.Background()
			onboardInfo, err := g.db.DeQOnboard(ctx)
			if err != nil {
				continue
			}

			if onboardInfo == nil {
				// g.log.Dbg("no pending employer onboard email generation")
				continue
			}

			g.log.Inf("onboard invites", "onboardInfo", onboardInfo)

			token := util.RandomUniqueID(vetchi.OrgUserInviteTokenLenBytes)
			link := vetchi.SignupOrgUserURL + token

			// TODO: Should migrate this to hedwig
			var textBody bytes.Buffer
			err = ttmpl.Must(
				ttmpl.New("text").Parse(textMailTemplate),
			).Execute(&textBody, map[string]string{
				"Domain": onboardInfo.DomainName,
				"Link":   link,
			})
			if err != nil {
				g.log.Err("email text template failed", "error", err)
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
				g.log.Err("email html template failed", "error", err)
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
