package hubopenings

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

func ApplyForOpening(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ApplyForOpening")
		var applyForOpeningReq hub.ApplyForOpeningRequest
		err := json.NewDecoder(r.Body).Decode(&applyForOpeningReq)
		if err != nil {
			h.Dbg("failed to decode apply for opening request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &applyForOpeningReq) {
			h.Dbg("validation failed")
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
				json.NewEncoder(w).Encode(common.ValidationErrors{
					Errors: []string{"resume"},
				})
				return
			}

			h.Dbg("failed to upload resume", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("uploaded resume", "filename", filename)

		applicationID := util.RandomUniqueID(vetchi.ApplicationIDLenBytes)
		h.Dbg("creating application in the db", "application_id", applicationID)

		err = h.DB().CreateApplication(r.Context(), db.ApplyOpeningReq{
			ApplicationID:          applicationID,
			OpeningIDWithinCompany: applyForOpeningReq.OpeningIDWithinCompany,
			CompanyDomain:          applyForOpeningReq.CompanyDomain,
			CoverLetter:            applyForOpeningReq.CoverLetter,
			ResumeSHA:              filename,
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
		err = json.NewEncoder(w).Encode(hub.ApplyForOpeningResponse{
			ApplicationID: applicationID,
		})
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

func uploadResume(
	ctx context.Context,
	h wand.Wand,
	resume string,
) (string, error) {
	// Validate and sanitize the PDF
	pdfBytes, err := util.ValidateAndSanitizePDF(resume)
	if err != nil {
		h.Dbg("invalid PDF file", "error", err)
		return "", db.ErrBadResume
	}

	_ = pdfBytes

	// S3TODO: Upload the resume to the object storage
	return "", nil
}
