package orgusers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func AddOrgUser(h wand.Wand) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddOrgUser")
		var addOrgUserReq employer.AddOrgUserRequest
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

		invite, err := generateOrgUserInvite(h, r, w, addOrgUserReq.Email)
		if err != nil {
			return
		}

		orgUserID, err := h.DB().AddOrgUser(r.Context(), db.AddOrgUserReq{
			Name:         addOrgUserReq.Name,
			Email:        addOrgUserReq.Email,
			OrgUserRoles: addOrgUserReq.Roles,
			OrgUserState: employer.AddedOrgUserState,

			InviteMail:  invite.Mail,
			InviteToken: invite.TokenReq,

			EmployerID:   invite.OrgUser.EmployerID,
			AddingUserID: invite.OrgUser.ID,
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
