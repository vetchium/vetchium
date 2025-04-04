package candidacy

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
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

func HubGetComments(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered HubGetComments")
		var getCommentsReq common.GetCandidacyCommentsRequest
		if err := json.NewDecoder(r.Body).Decode(&getCommentsReq); err != nil {
			h.Dbg("Error decoding request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getCommentsReq) {
			h.Dbg("Error validating request body")
			return
		}

		h.Dbg("validated", "getCommentsReq", getCommentsReq)

		comments, err := h.DB().
			GetHubCandidacyComments(r.Context(), getCommentsReq)
		if err != nil {
			h.Dbg("Error getting comments", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("got comments", "comments", comments)

		err = json.NewEncoder(w).Encode(comments)
		if err != nil {
			h.Dbg("Error encoding response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
