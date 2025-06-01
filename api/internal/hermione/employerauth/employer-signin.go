package employerauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"math/rand"
	"net/http"
	ttmpl "text/template"
	"time"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/employer"
	"golang.org/x/crypto/bcrypt"
)

func EmployerSignin(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered EmployerSignin")

		// Simulate a random delay to avoid timing attacks
		<-time.After(
			time.Millisecond * time.Duration(
				rand.Intn(int(h.Config().TimingAttackDelay.Milliseconds())),
			),
		)

		var employerSigninReq employer.EmployerSignInRequest

		err := json.NewDecoder(r.Body).Decode(&employerSigninReq)
		if err != nil {
			h.Dbg("failed to decode signin request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &employerSigninReq) {
			h.Dbg("validation failed")
			return
		}

		h.Dbg("validated",
			"clientID", employerSigninReq.ClientID,
			"email", employerSigninReq.Email,
		)

		orgUserAuth, err := h.DB().GetOrgUserAuth(
			r.Context(),
			db.OrgUserCreds{
				ClientID: employerSigninReq.ClientID,
				Email:    string(employerSigninReq.Email),
			},
		)
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				h.Dbg("no org user found",
					"clientID", employerSigninReq.ClientID,
					"email", employerSigninReq.Email,
				)
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			h.Dbg("database error getting org user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("org user found",
			"orgUserID", orgUserAuth.OrgUserID,
			"orgUserState", orgUserAuth.OrgUserState,
			"employerState", orgUserAuth.EmployerState,
		)

		if orgUserAuth.OrgUserState != employer.ActiveOrgUserState ||
			orgUserAuth.EmployerState != db.OnboardedEmployerState {
			h.Dbg("user or employer state check failed",
				"orgUserState", orgUserAuth.OrgUserState,
				"employerState", orgUserAuth.EmployerState,
			)
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(orgUserAuth.PasswordHash),
			[]byte(employerSigninReq.Password),
		)
		if err != nil {
			h.Dbg("password check failed", "error", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		h.Dbg("password check passed")

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
			h.Dbg("failed to init employer tfa", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(employer.EmployerSignInResponse{
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
Vetchium Team
`

	const htmlMailTemplate = `Hi,
<p>Please use the following token to signin to Vetchi.</p>
<p>Token: <b>{{.Token}}</b></p>
<p>Thanks,</p>
<p>Vetchium Team</p>
`

	const subject = "Vetchium Two Factor Authentication Token"

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
