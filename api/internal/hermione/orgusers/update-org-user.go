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

func UpdateOrgUser(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered UpdateOrgUser")

		var updateOrgUserReq vetchi.UpdateOrgUserRequest
		if err := json.NewDecoder(r.Body).Decode(&updateOrgUserReq); err != nil {
			h.Dbg("failed to decode update org user request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, updateOrgUserReq) {
			h.Dbg("validation failed", "updateOrgUserReq", updateOrgUserReq)
			return
		}

		if updateOrgUserReq.Name == "" && len(updateOrgUserReq.Roles) == 0 {
			http.Error(w, "name & roles cannot be empty", http.StatusBadRequest)
			return
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		orgUserID, err := h.DB().UpdateOrgUser(r.Context(), db.UpdateOrgUserReq{
			Name:           updateOrgUserReq.Name,
			Roles:          updateOrgUserReq.Roles,
			EmployerID:     orgUser.EmployerID,
			UpdatingUserID: orgUser.ID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				http.Error(w, err.Error(), http.StatusNotFound)
			}

			h.Dbg("failed to update org user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("org user updated", "orgUserID", orgUserID)
		w.WriteHeader(http.StatusOK)
	}
}
