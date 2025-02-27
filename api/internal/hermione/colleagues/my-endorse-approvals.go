package colleagues

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func MyEndorseApprovals(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered MyEndorseApprovals")
		var myEndorseApprovalsReq hub.MyEndorseApprovalsRequest
		err := json.NewDecoder(r.Body).Decode(&myEndorseApprovalsReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &myEndorseApprovalsReq) {
			h.Dbg("Validation failed")
			return
		}
		h.Dbg("validated", "myEndoreseApprovalsReq", myEndorseApprovalsReq)

		if myEndorseApprovalsReq.Limit == 0 {
			h.Dbg("Limit is 0, setting to default of 40")
			myEndorseApprovalsReq.Limit = 40
		}

		if myEndorseApprovalsReq.State == nil {
			h.Dbg("State is nil, setting to default of SoughtEndorsement")
			myEndorseApprovalsReq.State = []hub.EndorsementState{
				hub.SoughtEndorsement,
			}
		}

		endorsements, err := h.DB().
			GetMyEndorsementApprovals(r.Context(), myEndorseApprovalsReq)
		if err != nil {
			h.Dbg("Error getting endorsements", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Endorsements", "endorsements", endorsements)
		err = json.NewEncoder(w).Encode(endorsements)
		if err != nil {
			h.Dbg("Error encoding endorsements", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
