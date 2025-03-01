package profilepage

import (
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
)

// GetHubUserProfilePicture handles requests from employer to get a hub user's profile picture
func GetHubUserProfilePicture(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetHubUserProfilePicture")

		// Get the requested handle from URL path
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 { // /employer/get-hub-user-profile-picture/{handle}
			h.Dbg("invalid URL path", "path", r.URL.Path)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		requestedHandle := parts[4]

		// Get the profile picture URL for the requested handle
		pictureURL, err := h.DB().
			GetProfilePictureURL(r.Context(), requestedHandle)
		if err != nil {
			if err == db.ErrNoHubUser {
				h.Dbg("no hub user found", "handle", requestedHandle)
				http.Error(w, "", http.StatusNotFound)
				return
			}
			h.Dbg("failed to get profile picture URL", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// If no profile picture is set, return 404
		if pictureURL == "" {
			h.Dbg("no profile picture set", "handle", requestedHandle)
			http.Error(w, "", http.StatusNotFound)
			return
		}
		h.Dbg("got profile picture URL", "url", pictureURL)

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

		// Get the file from S3
		result, err := s3Client.GetObjectWithContext(
			r.Context(),
			&s3.GetObjectInput{
				Bucket: aws.String(cfg.S3.Bucket),
				Key:    aws.String(pictureURL),
			},
		)
		if err != nil {
			h.Err("failed to get object from S3", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer result.Body.Close()

		// Set content type based on file extension
		contentType := "image/jpeg" // default
		if strings.HasSuffix(pictureURL, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(pictureURL, ".webp") {
			contentType = "image/webp"
		}
		w.Header().Set("Content-Type", contentType)
		w.Header().
			Set("Cache-Control", "public, max-age=86400")
			// Cache for 24 hours

		h.Dbg("serving profile picture", "url", pictureURL)

		// Stream the file to the response
		_, err = io.Copy(w, result.Body)
		if err != nil {
			h.Err("failed to stream file", "error", err)
			// Can't write error response here as we've already started writing the response
		}
	}
}
