package granger

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vetchium/vetchium/api/internal/wand"
)

const maxFilesToCleanup = 100

// CleanupStaleFiles deletes old files from S3 and marks them as cleaned in the database
func CleanupStaleFiles(h wand.Wand) error {
	h.Dbg("starting cleanup of stale files")

	// Get unprocessed stale files
	staleFiles, err := h.DB().
		GetStaleFiles(context.Background(), maxFilesToCleanup)
	if err != nil {
		h.Err("failed to get stale files", "error", err)
		return fmt.Errorf("failed to get stale files: %w", err)
	}

	if len(staleFiles) == 0 {
		h.Dbg("no stale files to clean")
		return nil
	}

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

	// Process each stale file
	for _, file := range staleFiles {
		// Delete from S3
		_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(cfg.S3.Bucket),
			Key:    aws.String(file.FilePath),
		})
		if err != nil {
			h.Err("failed to delete file from S3",
				"error", err,
				"file_path", file.FilePath,
			)
			continue
		}

		// Mark as cleaned in database
		err = h.DB().MarkFileCleaned(
			context.Background(),
			file.ID,
			time.Now().UTC(),
		)
		if err != nil {
			h.Err("failed to mark file as cleaned",
				"error", err,
				"file_path", file.FilePath,
			)
			continue
		}

		h.Dbg("cleaned up stale file", "file_path", file.FilePath)
	}

	h.Dbg("completed cleanup of stale files",
		"processed_count", len(staleFiles),
	)
	return nil
}
