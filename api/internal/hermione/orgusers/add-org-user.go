package orgusers

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
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

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		orgUserID, err := h.DB().AddOrgUser(r.Context(), db.AddOrgUserReq{
			Email:        addOrgUserReq.Email,
			OrgUserRoles: addOrgUserReq.Roles,
			OrgUserState: db.AddedOrgUserState,
			EmployerID:   orgUser.EmployerID,
		})
		if err != nil {
			h.Dbg("AddOrgUser DB callfailed", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.Dbg("OrgUser Added", "orgUserID", orgUserID)

		w.WriteHeader(http.StatusOK)
	})
}
