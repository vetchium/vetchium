package locations

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/vhandler"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func DefunctLocation(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var defunctLocationRequest vetchi.DefunctLocationRequest
		err := json.NewDecoder(r.Body).Decode(&defunctLocationRequest)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &defunctLocationRequest) {
			return
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
		if !ok {
			h.Log().Error("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().DefunctLocation(r.Context(), db.DefunctLocationReq{
			Title:      defunctLocationRequest.Title,
			EmployerID: orgUser.EmployerID,
			OrgUserID:  orgUser.ID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoLocation) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
