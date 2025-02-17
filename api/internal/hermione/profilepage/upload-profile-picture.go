package profilepage

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var allowedFormats = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func UploadProfilePicture(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered UploadProfilePicture")

		// Parse multipart form with max memory of 5MB (our max file size limit)
		err := r.ParseMultipartForm(vetchi.MaxProfilePictureSize)
		if err != nil {
			h.Dbg("failed to parse multipart form", "error", err)
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		// Get the file from form
		file, header, err := r.FormFile("image")
		if err != nil {
			h.Dbg("failed to get file from form", "error", err)
			http.Error(w, "failed to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Validate image
		_, err = ValidateImage(
			file,
			header.Header.Get("Content-Type"),
			header.Size,
		)
		if err != nil {
			h.Dbg("image validation failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Reset file pointer after validation
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			h.Dbg("failed to reset file pointer", "error", err)
			http.Error(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
			return
		}

		// Generate unique filename with original extension
		ext := ".jpg" // default to jpg
		if strings.HasSuffix(strings.ToLower(header.Filename), ".png") {
			ext = ".png"
		} else if strings.HasSuffix(strings.ToLower(header.Filename), ".webp") {
			ext = ".webp"
		}

		// Generate unique ID for the file
		pictureID := util.RandomUniqueID(vetchi.ProfilePictureIDLenBytes)
		filename := fmt.Sprintf("profile-pictures/%s%s", pictureID, ext)
		h.Dbg(
			"generated filename",
			"filename",
			filename,
			"picture_id",
			pictureID,
		)

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
		_, err = s3Client.HeadBucketWithContext(
			r.Context(),
			&s3.HeadBucketInput{
				Bucket: aws.String(cfg.S3.Bucket),
			},
		)
		if err != nil {
			h.Dbg(
				"bucket does not exist, attempting to create",
				"bucket",
				cfg.S3.Bucket,
			)
			_, err = s3Client.CreateBucketWithContext(
				r.Context(),
				&s3.CreateBucketInput{
					Bucket: aws.String(cfg.S3.Bucket),
				},
			)
			if err != nil {
				h.Err("failed to create bucket", "error", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			h.Dbg("created bucket", "bucket", cfg.S3.Bucket)
		}

		// Upload to S3
		uploadInput := &s3.PutObjectInput{
			Bucket:        aws.String(cfg.S3.Bucket),
			Key:           aws.String(filename),
			Body:          file,
			ContentType:   aws.String(header.Header.Get("Content-Type")),
			ContentLength: aws.Int64(header.Size),
		}

		_, err = s3Client.PutObjectWithContext(r.Context(), uploadInput)
		if err != nil {
			h.Err("failed to upload profile picture", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// TODO: Implement UpdateProfilePicture method in the database interface
		// This method should update the user's profile picture filename in the database
		// The database should store both the pictureID and the full filename
		// err = h.DB().UpdateProfilePicture(r.Context(), pictureID, filename)
		// if err != nil {
		// 	h.Err("failed to update profile picture in database", "error", err)
		// 	http.Error(w, "", http.StatusInternalServerError)
		// 	return
		// }

		w.WriteHeader(http.StatusOK)
	}
}

// ValidateImage checks the file size, format, and dimensions
func ValidateImage(
	file multipart.File,
	contentType string,
	fileSize int64,
) (image.Image, error) {
	if fileSize > vetchi.MaxProfilePictureSize {
		return nil, fmt.Errorf("file exceeds max size of 5MB")
	}

	if !allowedFormats[contentType] {
		return nil, fmt.Errorf(
			"unsupported file format: only JPEG, PNG, and WEBP are allowed",
		)
	}

	// Read the image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Get dimensions
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width < vetchi.MinProfilePictureDim ||
		height < vetchi.MinProfilePictureDim ||
		width > vetchi.MaxProfilePictureDim ||
		height > vetchi.MaxProfilePictureDim {
		return nil, fmt.Errorf(
			"image dimensions must be between %dx%d and %dx%d",
			vetchi.MinProfilePictureDim,
			vetchi.MinProfilePictureDim,
			vetchi.MaxProfilePictureDim,
			vetchi.MaxProfilePictureDim,
		)
	}

	return img, nil
}
