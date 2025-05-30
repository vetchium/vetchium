package empposts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

func AddEmployerPost(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddEmployerPost")
		var addPostReq employer.AddEmployerPostRequest
		if err := json.NewDecoder(r.Body).Decode(&addPostReq); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addPostReq) {
			h.Dbg("validation failed")
			return
		}
		h.Dbg("validated", "req", addPostReq)

		postID := util.RandomUniqueID(vetchi.PostIDLenBytes)

		err := h.DB().AddEmployerPost(db.AddEmployerPostRequest{
			Context:                r.Context(),
			PostID:                 postID,
			AddEmployerPostRequest: addPostReq,
		})
		if err != nil {
			if errors.Is(err, db.ErrInvalidTagIDs) {
				h.Dbg("invalid tag IDs provided", "error", err)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(common.ValidationErrors{
					Errors: []string{"tags"},
				})
				if err != nil {
					h.Err("failed to encode validation errors", "error", err)
				}
				return
			}

			h.Dbg("failed to add employer post", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("employer post added", "post_id", postID)

		err = json.NewEncoder(w).Encode(employer.AddEmployerPostResponse{
			PostID: postID,
		})
		if err != nil {
			h.Err("encoding failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
