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
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"golang.org/x/crypto/bcrypt"
)

func EmployerSignin(h wand.Wand) http.HandlerFunc {
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

		if orgUserAuth.OrgUserState != vetchi.ActiveOrgUserState ||
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

		tfaMailCode, err := util.RandNumString(6)
		if err != nil {
			h.Dbg("failed to generate tfa mail code", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// We can even use the employerSigninReq.Email here but this
		// feels better. TODO: This needs to migrate to Hedwig package.
		email, err := generateEmail(orgUserAuth.OrgUserEmail, tfaMailCode)
		if err != nil {
			h.Dbg("failed to generate email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Ensures randomness and security
		tfaTokenString := util.RandomString(vetchi.TGTokenLenBytes)
		// Minimizes Collision and aspires for uniqueness
		tfaTokenString = tfaTokenString + time.Now().Format("0405")

		// TODO: Should we just email a magic URL instead of a token ? We can
		// make it longer, so minimize collisions and also more secure.

		err = h.DB().InitEmployerTFA(
			r.Context(),
			db.EmployerTFA{
				TFAToken: db.EmployerTokenReq{
					Token:            tfaTokenString,
					TokenType:        db.EmployerTFAToken,
					ValidityDuration: h.Config().Employer.TFATokLife,
					OrgUserID:        orgUserAuth.OrgUserID,
				},
				TFACode: tfaMailCode,
				Email:   email,
			},
		)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(vetchi.EmployerSignInResponse{
			Token: tfaTokenString,
		})
		if err != nil {
			h.Dbg("encode employer signin response", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// TODO: This needs to migrate to Hedwig package.
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
