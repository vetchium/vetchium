package orgusers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/vhandler"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func DisableOrgUser(h vhandler.VHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered DisableOrgUser")
		var disableOrgUserReq vetchi.DisableOrgUserRequest
		err := json.NewDecoder(r.Body).Decode(&disableOrgUserReq)
		if err != nil {
			h.Dbg("DisableOrgUserReq JSON decode failed", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Dbg("DisableOrgUserReq", "req", disableOrgUserReq)

		if !h.Vator().Struct(w, disableOrgUserReq) {
			h.Dbg("DisableOrgUserReq is not valid", "req", disableOrgUserReq)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().DisableOrgUser(r.Context(), db.DisableOrgUserReq{
			Email:           disableOrgUserReq.Email,
			EmployerID:      orgUser.EmployerID,
			DisablingUserID: orgUser.ID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				h.Dbg("DisableOrgUser: OrgUser not found", "err", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrLastActiveAdmin) {
				h.Dbg("DisableOrgUser: Last active admin", "err", err)
				http.Error(w, "last active admin", http.StatusForbidden)
				return
			}

			h.Err("DisableOrgUser DB call failed", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("orgUser disabled", "email", disableOrgUserReq.Email)
		w.WriteHeader(http.StatusOK)
	})
}
