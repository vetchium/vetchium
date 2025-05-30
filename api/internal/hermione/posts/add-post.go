package posts

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func AddPost(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddPost")
		var addPostReq hub.AddPostRequest
		if err := json.NewDecoder(r.Body).Decode(&addPostReq); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addPostReq) {
			h.Dbg("validation failed", "addPostReq", addPostReq)
			http.Error(w, "validation failed", http.StatusBadRequest)
			return
		}

		h.Dbg("Validated", "addPostReq", addPostReq)

		postID := util.RandomUniqueID(vetchi.PostIDLenBytes)

		err := h.DB().AddPost(db.AddPostRequest{
			Context:        r.Context(),
			PostID:         postID,
			AddPostRequest: addPostReq,
		})
		if err != nil {
			if err == db.ErrInvalidTagIDs {
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
			h.Dbg("adding post failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(hub.AddPostResponse{
			PostID: postID,
		})
		if err != nil {
			h.Dbg("encoding failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
