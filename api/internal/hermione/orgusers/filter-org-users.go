package orgusers

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func FilterOrgUsers(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered FilterOrgUsers")
		var filterOrgUsersReq employer.FilterOrgUsersRequest
		if err := json.NewDecoder(r.Body).Decode(&filterOrgUsersReq); err != nil {
			h.Err("failed to decode filter org users request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &filterOrgUsersReq) {
			h.Dbg("validation failed", "filterOrgUsersReq", filterOrgUsersReq)
			return
		}
		h.Dbg("validated", "filterOrgUsersReq", filterOrgUsersReq)

		if len(filterOrgUsersReq.State) == 0 {
			h.Dbg("no state specified, defaulting to ActiveOrgUserState")
			filterOrgUsersReq.State = []employer.OrgUserState{
				employer.ActiveOrgUserState,
				employer.AddedOrgUserState,
			}
		}

		if filterOrgUsersReq.Limit == 0 {
			h.Dbg("no limit specified, defaulting to 40")
			filterOrgUsersReq.Limit = 40
		}

		orgUsers, err := h.DB().
			FilterOrgUsers(r.Context(), filterOrgUsersReq)
		if err != nil {
			h.Dbg("failed to filter org users", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("filtered org users", "orgUsers", orgUsers)
		if err := json.NewEncoder(w).Encode(orgUsers); err != nil {
			h.Err("failed to encode org users", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
