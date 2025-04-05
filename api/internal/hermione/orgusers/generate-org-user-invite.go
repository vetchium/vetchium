package orgusers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
)

type invite struct {
	Mail     db.Email
	TokenReq db.OrgUserInviteReq
	OrgUser  db.OrgUserTO
}

func generateOrgUserInvite(
	h wand.Wand,
	r *http.Request,
	w http.ResponseWriter,
	affectedOrgUserEmail string,
) (invite, error) {
	orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		h.Err("failed to get orgUser from context")
		http.Error(w, "", http.StatusInternalServerError)
		return invite{}, errors.New("failed to get orgUser from context")
	}

	domains, err := h.DB().GetDomainNames(r.Context(), orgUser.EmployerID)
	if err != nil {
		h.Err("failed to get domains", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return invite{}, err
	}
	domainList := strings.Join(domains, ", ")

	// Ensures secrecy
	token := util.RandomString(vetchi.OrgUserInviteTokenLenBytes)
	// Ensures uniqueness. This is not needed mostly, but good to have
	token = token + strconv.FormatInt(time.Now().UnixNano(), 36)

	link := vetchi.EmployerBaseURL + "/signup-orguser/" + token

	inviteMail, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
		TemplateName: hedwig.InviteEmployee,
		Args: map[string]string{
			"Domains": domainList,
			"Link":    link,
		},
		EmailFrom: vetchi.EmailFrom,
		EmailTo:   []string{affectedOrgUserEmail},

		// TODO: The subject should be from Hedwig, based on the template
		// This subject is used in 0004-org-users_test.go too. Any change
		// in either place should be synced.
		Subject: "Vetchium Employer Invitation",
	})
	if err != nil {
		h.Dbg("failed to generate invite mail", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return invite{}, err
	}

	inviteTokenReq := db.OrgUserInviteReq{
		Token:            token,
		ValidityDuration: h.Config().Employer.InviteTokLife,
	}

	return invite{
		Mail:     inviteMail,
		TokenReq: inviteTokenReq,
		OrgUser:  orgUser,
	}, nil
}
