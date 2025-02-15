package applications

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func GetResume(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetResume")
		var getResumeRequest employer.GetResumeRequest
		err := json.NewDecoder(r.Body).Decode(&getResumeRequest)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Dbg("GetResume request", "request", getResumeRequest)

		if !h.Vator().Struct(w, &getResumeRequest) {
			h.Dbg("failed to validate request")
			return
		}
		h.Dbg("validated", "getResumeReq", getResumeRequest)

		// Get the resume details
		details, err := h.DB().
			GetResumeDetails(r.Context(), getResumeRequest)
		if err != nil {
			h.Dbg("failed to get resume details", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("got resume details", "details", details)

		// Create a descriptive filename
		filename := fmt.Sprintf(
			"%s-%s.pdf",
			details.ApplicationID,
			details.HubUserHandle,
		)
		h.Dbg("constructed filename", "filename", filename)

		// S3TODO: Read the resume from the object storage and serve it
	}
}
