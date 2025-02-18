package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
)

// UpdateProfilePictureWithCleanup updates a user's profile picture URL and adds the old
// picture to stale_files if one exists. This is done in a single transaction.
func (p *PG) UpdateProfilePictureWithCleanup(
	ctx context.Context,
	hubUserID uuid.UUID,
	newPicturePath string,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get current profile picture URL if exists
	var oldPicturePath *string
	err = tx.QueryRow(ctx, `
		SELECT profile_picture_url
		FROM hub_users
		WHERE id = $1
	`, hubUserID).Scan(&oldPicturePath)
	if err != nil {
		return fmt.Errorf("failed to get current profile picture: %w", err)
	}

	// Update profile picture URL
	result, err := tx.Exec(ctx, `
		UPDATE hub_users
		SET profile_picture_url = $1,
			updated_at = NOW()
		WHERE id = $2
	`, newPicturePath, hubUserID)
	if err != nil {
		return fmt.Errorf("failed to update profile picture: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user found with ID %s", hubUserID)
	}

	// If there was an old picture, add it to stale_files
	if oldPicturePath != nil && *oldPicturePath != "" {
		_, err = tx.Exec(ctx, `
			INSERT INTO stale_files (
				file_path
			) VALUES (
				$1
			)
		`, oldPicturePath)
		if err != nil {
			return fmt.Errorf(
				"failed to add old picture to stale files: %w",
				err,
			)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetStaleFiles retrieves a batch of unprocessed stale files
func (p *PG) GetStaleFiles(
	ctx context.Context,
	limit int,
) ([]db.StaleFile, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, file_path
		FROM stale_files
		WHERE cleaned_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query stale files: %w", err)
	}
	defer rows.Close()

	var files []db.StaleFile
	for rows.Next() {
		var file db.StaleFile
		err := rows.Scan(&file.ID, &file.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stale file: %w", err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stale files: %w", err)
	}

	return files, nil
}

// MarkFileCleaned marks a file as cleaned in the database
func (p *PG) MarkFileCleaned(
	ctx context.Context,
	fileID uuid.UUID,
	cleanedAt time.Time,
) error {
	result, err := p.pool.Exec(ctx, `
		UPDATE stale_files
		SET cleaned_at = $1
		WHERE id = $2
	`, cleanedAt, fileID)
	if err != nil {
		return fmt.Errorf("failed to mark file as cleaned: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no file found with ID %s", fileID)
	}

	return nil
}
