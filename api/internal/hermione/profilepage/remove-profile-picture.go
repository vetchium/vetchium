package profilepage

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/wand"
)

func RemoveProfilePicture(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RemoveProfilePicture")

		// Get hub user from context
		hubUser, ok := r.Context().Value(middleware.HubUserCtxKey).(db.HubUserTO)
		if !ok {
			h.Dbg("no hub user found in context")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// Get current profile picture URL if exists
		pictureURL, err := h.DB().
			GetProfilePictureURL(r.Context(), hubUser.Handle)
		if err != nil {
			if err == db.ErrNoHubUser {
				h.Dbg("no hub user found", "handle", hubUser.Handle)
				http.Error(w, "", http.StatusNotFound)
				return
			}
			h.Err("failed to get profile picture URL", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// If no profile picture is set, return 404
		if pictureURL == "" {
			h.Dbg("no profile picture set", "handle", hubUser.Handle)
			http.Error(w, "", http.StatusNotFound)
			return
		}

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

		// Delete the file from S3
		_, err = s3Client.DeleteObjectWithContext(
			r.Context(),
			&s3.DeleteObjectInput{
				Bucket: aws.String(cfg.S3.Bucket),
				Key:    aws.String(pictureURL),
			},
		)
		if err != nil {
			h.Err("failed to delete object from S3", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Clear the profile picture URL in the database
		err = h.DB().
			UpdateProfilePictureWithCleanup(r.Context(), hubUser.ID, "")
		if err != nil {
			h.Err("failed to clear profile picture URL", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("removed profile picture", "handle", hubUser.Handle)
		w.WriteHeader(http.StatusOK)
	}
}
