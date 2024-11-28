package hubopenings

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func ApplyForOpeningHandler(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ApplyForOpeningHandler")
		var applyForOpeningReq vetchi.ApplyForOpeningRequest
		err := json.NewDecoder(r.Body).Decode(&applyForOpeningReq)
		if err != nil {
			h.Dbg("failed to decode apply for opening request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &applyForOpeningReq) {
			h.Dbg("validation failed", "applyForOpeningReq", applyForOpeningReq)
			return
		}
		h.Dbg("validated", "applyForOpeningReq", applyForOpeningReq)

		resumeFileReq, err := http.NewRequestWithContext(
			r.Context(),
			http.MethodPost,
			"http://granger:8080/internal/upload-resume", // Take this from config ?
			bytes.NewBuffer([]byte(applyForOpeningReq.Resume)),
		)
		if err != nil {
			h.Err("failed to create resume file request", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		resumeFileResp, err := http.DefaultClient.Do(resumeFileReq)
		if err != nil {
			h.Err("failed to upload resume", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		defer resumeFileResp.Body.Close()

		if resumeFileResp.StatusCode != http.StatusOK {
			h.Err("failed to upload resume", "status", resumeFileResp.Status)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		resumeFileRespBody, err := io.ReadAll(resumeFileResp.Body)
		if err != nil {
			h.Err("failed to read resume file response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		filename := string(resumeFileRespBody)
		h.Dbg("uploaded resume", "filename", filename)

		// Ensures secrecy
		applicationID := util.RandomString(vetchi.ApplicationIDLenBytes)
		// Ensures uniqueness. This is not needed mostly, but good to have
		applicationID = applicationID + strconv.FormatInt(
			time.Now().UnixNano(),
			36,
		)

		err = h.DB().CreateApplication(r.Context(), db.ApplyOpeningReq{
			ApplicationID:          applicationID,
			OpeningIDWithinCompany: applyForOpeningReq.OpeningIDWithinCompany,
			CompanyDomain:          applyForOpeningReq.CompanyDomain,
			CoverLetter:            applyForOpeningReq.CoverLetter,
			OriginalFilename:       applyForOpeningReq.Filename,
			InternalFilename:       filename,
		})
		if err != nil {
			h.Err("failed to create application", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("created application", "application_id", applicationID)
		w.WriteHeader(http.StatusOK)
	}
}
