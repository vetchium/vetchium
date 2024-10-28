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

func RenameLocation(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var renameLocationReq vetchi.RenameLocationRequest
		err := json.NewDecoder(r.Body).Decode(&renameLocationReq)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &renameLocationReq) {
			return
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
		if !ok {
			h.Log().Error("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().RenameLocation(r.Context(), db.RenameLocationReq{
			EmployerID: orgUser.EmployerID,
			OldTitle:   renameLocationReq.OldTitle,
			NewTitle:   renameLocationReq.NewTitle,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoLocation) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrDupLocationName) {
				http.Error(w, "", http.StatusConflict)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
