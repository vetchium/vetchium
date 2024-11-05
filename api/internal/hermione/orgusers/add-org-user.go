package orgusers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func AddOrgUser(h wand.Wand) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddOrgUser")
		var addOrgUserReq vetchi.AddOrgUserRequest
		err := json.NewDecoder(r.Body).Decode(&addOrgUserReq)
		if err != nil {
			h.Dbg("AddOrgUserReq JSON decode failed", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Dbg("AddOrgUserReq", "req", addOrgUserReq)

		if !h.Vator().Struct(w, &addOrgUserReq) {
			h.Dbg("AddOrgUserReq validation failed", "req", addOrgUserReq)
			return
		}
		h.Dbg("validated", "AddOrgUserReq", addOrgUserReq)

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		domains, err := h.DB().GetDomainNames(r.Context(), orgUser.EmployerID)
		if err != nil {
			h.Dbg("failed to get domains", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		domainList := strings.Join(domains, ", ")

		// Ensures secrecy
		token := util.RandomString(vetchi.OnBoardTokenLenBytes)
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
			EmailTo:   []string{addOrgUserReq.Email},

			// TODO: The subject should be from Hedwig, based on the template
			// This subject is used in 0004-org-users_test.go too. Any change
			// in either place should be synced.
			Subject: "Vetchi Employer Invitation",
		})
		if err != nil {
			h.Dbg("failed to generate invite mail", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		inviteTokenValidityDuration, err := h.ConfigDuration(
			db.EmployerInviteToken,
		)
		if err != nil {
			h.Dbg("failed to get invite token validity duration", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		orgUserID, err := h.DB().AddOrgUser(r.Context(), db.AddOrgUserReq{
			Name:         addOrgUserReq.Name,
			Email:        addOrgUserReq.Email,
			OrgUserRoles: addOrgUserReq.Roles,
			OrgUserState: vetchi.AddedOrgUserState,

			InviteMail: inviteMail,
			InviteToken: db.TokenReq{
				Token:            token,
				TokenType:        db.EmployerInviteToken,
				ValidityDuration: inviteTokenValidityDuration,
			},

			EmployerID:   orgUser.EmployerID,
			AddingUserID: orgUser.ID,
		})
		if err != nil {
			if errors.Is(err, db.ErrOrgUserAlreadyExists) {
				h.Dbg("org user already exists", "addOrgUserReq", addOrgUserReq)
				http.Error(w, "", http.StatusConflict)
				return
			}

			h.Dbg("AddOrgUser DB callfailed", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("OrgUser Added", "orgUserID", orgUserID)

		w.WriteHeader(http.StatusOK)
	})
}
