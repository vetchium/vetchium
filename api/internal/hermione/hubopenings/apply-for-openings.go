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

	// Calculate SHA-512 hash of the PDF content
	hash := sha512.Sum512(pdfBytes)
	filename := fmt.Sprintf("%x.pdf", hash)
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
