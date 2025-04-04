package applications

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
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
		details, err := h.DB().GetResumeDetails(r.Context(), getResumeRequest)
		if err != nil {
			h.Dbg("failed to get resume details", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("got resume details", "details", details)

		// Create a descriptive filename for download
		filename := fmt.Sprintf(
			"%s-%s.pdf",
			details.ApplicationID,
			details.HubUserHandle,
		)
		h.Dbg("constructed filename", "filename", filename)

		// Initialize S3 client
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
		s3Client := s3.New(session.Must(session.NewSession(s3Config)))

		// Get the file from S3
		result, err := s3Client.GetObjectWithContext(
			r.Context(),
			&s3.GetObjectInput{
				Bucket: aws.String(cfg.S3.Bucket),
				Key:    aws.String(details.SHA),
			},
		)
		if err != nil {
			h.Err("failed to get resume from S3", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer result.Body.Close()

		// Set response headers for PDF viewing
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().
			Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
		if result.ContentLength != nil {
			w.Header().
				Set("Content-Length", fmt.Sprintf("%d", *result.ContentLength))
		}

		// Stream the file to the response
		_, err = io.Copy(w, result.Body)
		if err != nil {
			h.Err("failed to stream resume to response", "error", err)
			// Note: Headers might have been sent already, so we can't send an error response
			return
		}
		h.Dbg("successfully served resume", "filename", filename)
	}
}
