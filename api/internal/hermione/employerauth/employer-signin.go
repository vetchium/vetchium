package employerauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	ttmpl "text/template"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/vhandler"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"golang.org/x/crypto/bcrypt"
)

func EmployerSignin(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employerSigninReq vetchi.EmployerSignInRequest

		err := json.NewDecoder(r.Body).Decode(&employerSigninReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &employerSigninReq) {
			return
		}

		orgUserAuth, err := h.DB().GetOrgUserAuth(
			r.Context(),
			db.OrgUserCreds{
				ClientID: employerSigninReq.ClientID,
				Email:    string(employerSigninReq.Email),
			},
		)
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if orgUserAuth.OrgUserState != db.ActiveOrgUserState ||
			orgUserAuth.EmployerState != db.OnboardedEmployerState {
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(orgUserAuth.PasswordHash),
			[]byte(employerSigninReq.Password),
		)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		tgToken := db.OrgUserToken{
			Token:          util.RandomString(vetchi.TGTokenLenBytes),
			OrgUserID:      orgUserAuth.OrgUserID,
			TokenValidTill: time.Now().Add(h.TGTLife()),
			TokenType:      db.TGToken,
		}

		emailTokenString := util.RandomString(vetchi.EmailTokenLenBytes)
		emailToken := db.OrgUserToken{
			Token:          emailTokenString,
			OrgUserID:      orgUserAuth.OrgUserID,
			TokenValidTill: time.Now().Add(h.TGTLife()),
			TokenType:      db.EmailToken,
		}

		// We can even use the employerSigninReq.Email here but this
		// feels better.
		email, err := generateEmail(orgUserAuth.OrgUserEmail, emailTokenString)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().InitEmployerTFA(
			r.Context(),
			db.EmployerTFA{
				TGToken:    tgToken,
				EmailToken: emailToken,
				Email:      email,
			},
		)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(vetchi.EmployerSignInResponse{
			Token: tgToken.Token,
		})
		if err != nil {
			h.Err(
				"failed to encode employer signin response",
				"error",
				err,
			)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func generateEmail(orgUserEmail, token string) (db.Email, error) {
	const textMailTemplate = `
Hi there,

Please use the following token to signin to Vetchi.

Token: {{.Token}}

Thanks,
Vetchi Team
`

	const htmlMailTemplate = `Hi,
<p>Please use the following token to signin to Vetchi.</p>
<p>Token: <b>{{.Token}}</b></p>
<p>Thanks,</p>
<p>Vetchi Team</p>
`

	const subject = "Vetchi Two Factor Authentication Token"

	var textBody bytes.Buffer
	err := ttmpl.Must(
		ttmpl.New("text").Parse(textMailTemplate),
	).Execute(&textBody, map[string]string{
		"Token": token,
	})
	if err != nil {
		return db.Email{}, err
	}

	var htmlBody bytes.Buffer
	err = template.Must(
		template.New("html").Parse(htmlMailTemplate),
	).Execute(&htmlBody, map[string]string{
		"Token": token,
	})
	if err != nil {
		return db.Email{}, err
	}

	email := db.Email{
		EmailFrom:     vetchi.EmailFrom,
		EmailTo:       []string{orgUserEmail},
		EmailSubject:  subject,
		EmailHTMLBody: htmlBody.String(),
		EmailTextBody: textBody.String(),
		EmailState:    db.EmailStatePending,
	}

	return email, nil
}
