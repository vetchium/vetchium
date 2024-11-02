package orgusers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/vhandler"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func AddOrgUser(h vhandler.VHandler) http.HandlerFunc {
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

		inviteMail, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: hedwig.InviteEmployee,
			Args: map[string]string{
				"Domains": domainList,
			},
		})
		if err != nil {
			h.Dbg("failed to generate invite mail", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		orgUserID, err := h.DB().AddOrgUser(r.Context(), db.AddOrgUserReq{
			Name:         addOrgUserReq.Name,
			Email:        addOrgUserReq.Email,
			OrgUserRoles: addOrgUserReq.Roles,
			OrgUserState: vetchi.AddedOrgUserState,

			InviteMail: inviteMail,

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
