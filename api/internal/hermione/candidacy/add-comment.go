package candidacy

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func EmployerAddComment(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered EmployerAddComment")
		var addCommentReq employer.AddEmployerCandidacyCommentRequest
		err := json.NewDecoder(r.Body).Decode(&addCommentReq)
		if err != nil {
			h.Dbg("Error decoding request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addCommentReq) {
			h.Dbg("Validation failed")
			return
		}
		h.Dbg("validated", "addCommentReq", addCommentReq)

		commentID, err := h.DB().
			AddEmployerCandidacyComment(r.Context(), addCommentReq)
		if err != nil {
			switch err {
			case db.ErrNoOpening:
				h.Dbg("Candidacy not found", "error", err)
				http.Error(w, "Candidacy not found", http.StatusNotFound)
			case db.ErrInvalidCandidacyState:
				h.Dbg("Invalid candidacy state", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
			case db.ErrUnauthorizedComment:
				h.Dbg("User not authorized to comment", "error", err)
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
