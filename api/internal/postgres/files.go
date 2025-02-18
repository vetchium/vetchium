package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/psankar/vetchi/api/internal/db"
)

// UpdateProfilePictureWithCleanup updates a user's profile picture URL and adds the old
// picture to stale_files if one exists. This is done in a single transaction.
// If newPicturePath is empty, the profile_picture_url will be set to NULL.
func (p *PG) UpdateProfilePictureWithCleanup(
	ctx context.Context,
	hubUserID uuid.UUID,
	newPicturePath string,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
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
		p.log.Err("failed to get current profile picture", "error", err)
		return fmt.Errorf("failed to get current profile picture: %w", err)
	}

	// Update profile picture URL
	// If newPicturePath is empty, set to NULL
	var result pgconn.CommandTag
	if newPicturePath == "" {
		p.log.Dbg("updating profile picture to NULL", "user_id", hubUserID)
		result, err = tx.Exec(ctx, `
			UPDATE hub_users
			SET profile_picture_url = NULL,
				updated_at = NOW()
			WHERE id = $1
		`, hubUserID)
	} else {
		p.log.Dbg("updating path", "user_id", hubUserID, "path", newPicturePath)
		result, err = tx.Exec(ctx, `
			UPDATE hub_users
			SET profile_picture_url = $1,
				updated_at = NOW()
			WHERE id = $2
		`, newPicturePath, hubUserID)
	}

	if err != nil {
		p.log.Dbg("failed to update profile picture", "error", err)
		return fmt.Errorf("failed to update profile picture: %w", err)
	}

	if result.RowsAffected() == 0 {
		p.log.Dbg("no user found with ID", "id", hubUserID)
		return fmt.Errorf("no user found with ID %s", hubUserID)
	}

	// If there was an old picture, add it to stale_files
	if oldPicturePath != nil && *oldPicturePath != "" {
		p.log.Dbg("adding old picture to stale files", "path", *oldPicturePath)
		_, err = tx.Exec(ctx, `
			INSERT INTO stale_files (
				file_path
			) VALUES (
				$1
			)
		`, oldPicturePath)
		if err != nil {
			p.log.Err("failed to add old picture to stale files", "error", err)
			return fmt.Errorf(
				"failed to add old picture to stale files: %w",
				err,
			)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	p.log.Dbg("updated profile picture", "user_id", hubUserID)
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

// GetProfilePictureURL returns the S3 URL of a user's profile picture
func (p *PG) GetProfilePictureURL(
	ctx context.Context,
	handle string,
) (string, error) {
	var pictureURL *string
	err := p.pool.QueryRow(ctx, `
		SELECT profile_picture_url
		FROM hub_users
		WHERE handle = $1
	`, handle).Scan(&pictureURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", db.ErrNoHubUser
		}
		return "", fmt.Errorf("failed to get profile picture URL: %w", err)
	}
	if pictureURL == nil {
		return "", nil
	}
	return *pictureURL, nil
}
