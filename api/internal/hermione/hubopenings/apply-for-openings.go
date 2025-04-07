package hubopenings

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
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

		// TODO: Fetch the opening's title from the db
		openingTitle := ""
		openingURL := h.Config().Hub.WebURL + "/org/" +
			applyForOpeningReq.CompanyDomain + "/opening/" +
			applyForOpeningReq.OpeningIDWithinCompany

		// If there are endorsers, prepare endorsement emails
		var endorsementEmails []db.Email
		if len(applyForOpeningReq.EndorserHandles) > 0 {
			endorsementEmails, err = prepareEndorsementEmails(
				r.Context(),
				h,
				w,
				applyForOpeningReq.EndorserHandles,
				applyForOpeningReq.CompanyDomain,
				openingTitle,
				openingURL,
			)
			if err != nil {
				h.Dbg("failed to prepare endorsement emails", "error", err)
				return
			}
		}

		err = h.DB().CreateApplication(r.Context(), db.ApplyOpeningReq{
			ApplicationID:          applicationID,
			OpeningIDWithinCompany: applyForOpeningReq.OpeningIDWithinCompany,
			CompanyDomain:          applyForOpeningReq.CompanyDomain,
			CoverLetter:            applyForOpeningReq.CoverLetter,
			ResumeSHA:              filename,
			EndorserHandles:        applyForOpeningReq.EndorserHandles,
			EndorsementEmails:      endorsementEmails,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("either domain or opening does not exist", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			// This is unlikely to happen, but handling gracefully if it does
			if errors.Is(err, db.ErrNotColleague) {
				h.Dbg("one or more endorsers are not colleagues", "error", err)
				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(common.ValidationErrors{
					Errors: []string{"endorser_handles"},
				})
				return
			}

			if errors.Is(err, db.ErrCannotApply) {
				h.Dbg("user cannot apply to this opening", "error", err)
				w.WriteHeader(http.StatusUnprocessableEntity)
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

	// Calculate SHA-512 hash of the PDF content
	hash := sha512.Sum512(pdfBytes)
	filename := fmt.Sprintf("%s%x.pdf", util.ResumesPath, hash)
	h.Dbg("calculated file hash", "sha512", filename)

	cfg := h.Config()
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			cfg.S3.AccessKey,
			cfg.S3.SecretKey,
			"",
		),
		Endpoint:         aws.String(cfg.S3.Endpoint),
		Region:           aws.String(cfg.S3.Region),
		S3ForcePathStyle: aws.Bool(true), // Required for MinIO
	}

	// Create S3 service client
	s3Client := s3.New(session.Must(session.NewSession(s3Config)))

	// Ensure bucket exists
	_, err = s3Client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.S3.Bucket),
	})
	if err != nil {
		h.Dbg(
			"bucket does not exist, attempting to create",
			"bucket",
			cfg.S3.Bucket,
		)
		_, err = s3Client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(cfg.S3.Bucket),
		})
		if err != nil {
			h.Err("failed to create bucket", "error", err)
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
		h.Dbg("created bucket", "bucket", cfg.S3.Bucket)
	}

	// Check if file already exists
	_, err = s3Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(cfg.S3.Bucket),
		Key:    aws.String(filename),
	})
	if err == nil {
		// File already exists, no need to upload again
		h.Dbg(
			"resume already exists in storage, skipping upload",
			"filename",
			filename,
		)
		return filename, nil
	}

	// File doesn't exist, proceed with upload
	uploadInput := &s3.PutObjectInput{
		Bucket:        aws.String(cfg.S3.Bucket),
		Key:           aws.String(filename),
		Body:          bytes.NewReader(pdfBytes),
		ContentType:   aws.String("application/pdf"),
		ContentLength: aws.Int64(int64(len(pdfBytes))),
	}

	// Upload the file
	_, err = s3Client.PutObjectWithContext(ctx, uploadInput)
	if err != nil {
		h.Err("failed to upload resume to S3", "error", err)
		return "", fmt.Errorf("failed to upload resume: %w", err)
	}

	h.Dbg("uploaded new resume", "filename", filename)
	return filename, nil
}

func prepareEndorsementEmails(
	ctx context.Context,
	h wand.Wand,
	w http.ResponseWriter,
	endorserHandles []common.Handle,
	companyName, openingTitle, openingURL string,
) ([]db.Email, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		h.Err("failed to get hub user from context")
		http.Error(w, "", http.StatusInternalServerError)
		return []db.Email{}, errors.New("failed to get hub user from context")
	}

	endorsementEmails := []db.Email{}
	// Get all endorser details in one call
	endorsers, err := h.DB().GetHubUsersByHandles(ctx, endorserHandles)
	if err != nil {
		h.Dbg("failed to get endorser details", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return []db.Email{}, errors.New("failed to get endorser details")
	}

	if len(endorsers) != len(endorserHandles) {
		h.Dbg("Duplicate/Missing endorsers", "endorserHandles", endorserHandles)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.ValidationErrors{
			Errors: []string{"endorser_handles"},
		})
		return []db.Email{}, errors.New("duplicate/missing endorsers")
	}

	// Create an email for each endorser
	for _, endorser := range endorsers {
		email, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: hedwig.EndorsementRequest,
			Args: map[string]string{
				"endorser_name":    endorser.FullName,
				"applicant_name":   hubUser.FullName,
				"applicant_handle": hubUser.Handle,
				"company_name":     companyName,
				"job_title":        openingTitle,
				"endorse_url":      h.Config().Hub.WebURL + "/my-approvals",
				"opening_url":      openingURL,
			},
			EmailFrom: vetchi.EmailFrom,
			EmailTo:   []string{endorser.Email},

			// TODO: This should be i18n enabled and come from the mail templates.
			Subject: "Endorsement Request",
		})
		if err != nil {
			h.Dbg("failed to generate endorsement email", "error", err)
			return []db.Email{}, err
		}

		endorsementEmails = append(endorsementEmails, email)
	}

	return endorsementEmails, nil
}
