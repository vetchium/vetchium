package applications

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/psankar/vetchi/api/internal/util"
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

		// Construct the file path using the SHA
		resumePath := util.GetResumeStoragePath("/resumes", details.SHA)
		h.Dbg("constructed resume path", "path", resumePath)

		// Open and serve the file
		file, err := os.Open(resumePath)
		if err != nil {
			h.Dbg("failed to open resume file", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Get file info for ModTime
		fileInfo, err := file.Stat()
		if err != nil {
			h.Dbg("failed to get file info", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Create a descriptive filename
		filename := fmt.Sprintf(
			"%s-%s.pdf",
			details.ApplicationID,
			details.HubUserHandle,
		)

		// Set content type to application/pdf
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().
			Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

		http.ServeContent(w, r, filename, fileInfo.ModTime(), file)
	}
}
