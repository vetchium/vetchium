package util

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
)

// AllowedImageFormats defines the allowed image MIME types
var AllowedImageFormats = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

const (
	// Profile picture specific constants
	MaxProfilePictureSize    = 5 * 1024 * 1024 // 5MB
	MinProfilePictureDim     = 200             // 200x200 pixels
	MaxProfilePictureDim     = 2048            // 2048x2048 pixels
	ProfilePictureIDLenBytes = 16              // Length of the unique ID for profile pictures

	// S3 storage paths
	ProfilePicturesPath = "hub-users/profile-pictures/" // Scoped under hub-users since it's user specific
	ResumesPath         = "resumes/"                    // Top-level since resumes can come from multiple sources
)

// ValidateImage checks if the given image file meets the size, format, and dimension requirements
// Returns the decoded image and error if any validation fails
func ValidateImage(
	file multipart.File,
	contentType string,
	fileSize int64,
	maxSize int64,
	minDim int,
	maxDim int,
) (image.Image, error) {
	if fileSize > maxSize {
		return nil, fmt.Errorf("file exceeds max size of %d bytes", maxSize)
	}

	if !AllowedImageFormats[contentType] {
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
	if width < minDim ||
		height < minDim ||
		width > maxDim ||
		height > maxDim {
		return nil, fmt.Errorf(
			"image dimensions must be between %dx%d and %dx%d",
			minDim,
			minDim,
			maxDim,
			maxDim,
		)
	}

	return img, nil
}

// ValidateProfilePicture is a convenience function that validates an image
// using the profile picture specific constraints
func ValidateProfilePicture(
	file multipart.File,
	contentType string,
	fileSize int64,
) (image.Image, error) {
	return ValidateImage(
		file,
		contentType,
		fileSize,
		MaxProfilePictureSize,
		MinProfilePictureDim,
		MaxProfilePictureDim,
	)
}
