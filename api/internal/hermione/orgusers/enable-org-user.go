package orgusers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func EnableOrgUser(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered EnableOrgUser")
		var enableOrgUserReq vetchi.EnableOrgUserRequest
		err := json.NewDecoder(r.Body).Decode(&enableOrgUserReq)
		if err != nil {
			h.Dbg("EnableOrgUserReq JSON decode failed", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Dbg("EnableOrgUserReq", "enableOrgUserReq", enableOrgUserReq)

		if !h.Vator().Struct(w, &enableOrgUserReq) {
			h.Dbg("EnableOrgUserReq is not valid", "req", enableOrgUserReq)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		h.Dbg("validated", "enableOrgUserReq", enableOrgUserReq)

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().EnableOrgUser(r.Context(), db.EnableOrgUserReq{
			Email:          enableOrgUserReq.Email,
			EmployerID:     orgUser.EmployerID,
			EnablingUserID: orgUser.ID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				h.Dbg("org user not found", "err", err)
				http.Error(w, "Org user not found", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrOrgUserNotDisabled) {
				h.Dbg("org user not disabled", "err", err)
				http.Error(w, "Org user not disabled", http.StatusBadRequest)
				return
			}

			h.Dbg("failed to enable org user", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("enabled org user", "email", enableOrgUserReq.Email)
		w.WriteHeader(http.StatusOK)
	}
}
