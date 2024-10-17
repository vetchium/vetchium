package hermione

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"golang.org/x/crypto/bcrypt"
)

func (h *Hermione) employerSignin(w http.ResponseWriter, r *http.Request) {
	var employerSigninReq vetchi.EmployerSignInRequest

	err := json.NewDecoder(r.Body).Decode(&employerSigninReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.vator.Struct(w, employerSigninReq) {
		return
	}

	orgUserAuth, err := h.db.GetOrgUserAuth(
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

	sessionToken := util.RandomString(vetchi.SessionTokenLenBytes)
	err = h.db.CreateOrgUserSession(
		r.Context(),
		db.OrgUserSession{
			OrgUserID:           orgUserAuth.OrgUserID,
			SessionToken:        sessionToken,
			SessionValidityMins: h.employer.tgtLife.Minutes(),
			TokenType:           db.TGToken,
		},
	)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	employerSigninResp := vetchi.EmployerSignInResponse{
		Token: sessionToken,
	}

	err = json.NewEncoder(w).Encode(employerSigninResp)
	if err != nil {
		h.log.Error("failed to encode employer signin response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
