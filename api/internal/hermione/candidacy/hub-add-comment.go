package candidacy

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func HubAddComment(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered HubAddComment")
		var addCommentReq hub.AddHubCandidacyCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&addCommentReq); err != nil {
			h.Dbg("Error decoding request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addCommentReq) {
			h.Dbg("Error validating request body")
			return
		}
		h.Dbg("validated", "addCommentReq", addCommentReq)

		commentID, err := h.DB().
			AddHubCandidacyComment(r.Context(), addCommentReq)
		if err != nil {
			h.Dbg("Error adding comment", "error", err)
			switch err {
			case db.ErrNoOpening:
				h.Dbg("Candidacy not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
			case db.ErrInvalidCandidacyState:
				h.Dbg("Invalid candidacy state", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
			case db.ErrUnauthorizedComment:
				h.Dbg("Unauthorized to comment", "error", err)
				http.Error(w, "", http.StatusForbidden)
			default:
				h.Dbg("Internal error while adding comment", "error", err)
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		h.Dbg("Added comment", "commentID", commentID)

		w.WriteHeader(http.StatusOK)
	}
}
