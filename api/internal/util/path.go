package util

import (
	"crypto/sha512"
	"encoding/hex"
	"path/filepath"
)

// GenerateResumeID generates a unique ID for a resume based on its content
// It returns a hex-encoded SHA-512 hash of the content
func GenerateResumeID(content []byte) string {
	hash := sha512.Sum512(content)
	return hex.EncodeToString(hash[:])
}

// GetResumeStoragePath returns the hierarchical path for storing a resume file
// It creates a directory structure using the first 6 characters of the SHA-512 hash
// in groups of 2, giving us 16^2 = 256 possibilities for each level
// With 3 levels, we get 256^3 = ~16.7 million possible buckets
// Example: for hash "a1b2c3d4..." returns "/resumes/a1/b2/c3/a1b2c3d4..."
func GetResumeStoragePath(baseDir string, resumeID string) string {
	if len(resumeID) < 6 {
		return filepath.Join(baseDir, resumeID)
	}

	dir1 := resumeID[0:2] // First 2 chars
	dir2 := resumeID[2:4] // Next 2 chars
	dir3 := resumeID[4:6] // Next 2 chars

	return filepath.Join(baseDir, dir1, dir2, dir3, resumeID)
}

// GetResumeStorageDir returns just the directory path where the resume should be stored
// Example: for hash "a1b2c3d4..." returns "/resumes/a1/b2/c3"
func GetResumeStorageDir(baseDir string, resumeID string) string {
	if len(resumeID) < 6 {
		return baseDir
	}

	dir1 := resumeID[0:2] // First 2 chars
	dir2 := resumeID[2:4] // Next 2 chars
	dir3 := resumeID[4:6] // Next 2 chars

	return filepath.Join(baseDir, dir1, dir2, dir3)
}
