package profilepage

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetMyDetails(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hubUser, ok := r.Context().Value(middleware.HubUserCtxKey).(db.HubUserTO)
		if !ok {
			h.Err("failed to get hub user")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		var myDetails hub.MyDetails
		myDetails.Handle = hubUser.Handle
		myDetails.FullName = hubUser.FullName
		myDetails.Tier = hubUser.Tier

		h.Dbg("my details", "myDetails", myDetails)

		if err := json.NewEncoder(w).Encode(myDetails); err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
