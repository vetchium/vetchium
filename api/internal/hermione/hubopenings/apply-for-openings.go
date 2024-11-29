package hubopenings

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

		// TODO: Validate if this hubUser can apply for this opening
		// Some essential but not complete list of things to check:
		// - Has the HubUser already applied for this Opening ?
		// - Has the HubUser already applied to this Employer in the last X months ?
		// - Is this an internal opening for the Employer ?
		// - Has the Employer blocked this HubUser ?
		// - Should we cross check against the Opening's Years of Experience expectations ?

		filename, err := uploadResume(r.Context(), h, applyForOpeningReq.Resume)
		if err != nil {
			if errors.Is(err, db.ErrBadResume) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(vetchi.ValidationErrors{
					Errors: []string{"resume"},
				})
				return
			}

			h.Dbg("failed to upload resume", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("uploaded resume", "filename", filename)

		// Ensures secrecy
		applicationID := util.RandomString(vetchi.ApplicationIDLenBytes)
		// Ensures uniqueness. This is not needed mostly, but good to have
		applicationID = applicationID + strconv.FormatInt(
			time.Now().UnixNano(),
			36,
		)

		h.Dbg("creating application in the db", "application_id", applicationID)

		err = h.DB().CreateApplication(r.Context(), db.ApplyOpeningReq{
			ApplicationID:          applicationID,
			OpeningIDWithinCompany: applyForOpeningReq.OpeningIDWithinCompany,
			CompanyDomain:          applyForOpeningReq.CompanyDomain,
			CoverLetter:            applyForOpeningReq.CoverLetter,
			OriginalFilename:       applyForOpeningReq.Filename,
			InternalFilename:       filename,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("either domain or opening does not exist", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Err("failed to create application", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("created application", "application_id", applicationID)
		w.WriteHeader(http.StatusOK)
	}
}

func uploadResume(
	ctx context.Context,
	h wand.Wand,
	resume string,
) (string, error) {
	resumeFileReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://granger:8080/internal/upload-resume", // Take this from config ?
		bytes.NewBuffer([]byte(resume)),
	)
	if err != nil {
		h.Err("failed to create resume file request", "error", err)
		return "", db.ErrInternal
	}

	resumeFileResp, err := http.DefaultClient.Do(resumeFileReq)
	if err != nil {
		h.Err("failed to upload resume", "error", err)
		return "", db.ErrInternal
	}

	defer resumeFileResp.Body.Close()

	switch resumeFileResp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		h.Dbg("failed to upload resume", "status", resumeFileResp.Status)
		return "", db.ErrBadResume
	default:
		h.Err("failed to upload resume", "status", resumeFileResp.Status)
		return "", db.ErrInternal
	}

	resumeFileRespBody, err := io.ReadAll(resumeFileResp.Body)
	if err != nil {
		h.Err("failed to read resume file response", "error", err)
		return "", db.ErrInternal
	}

	filename := string(resumeFileRespBody)
	return filename, nil
}
