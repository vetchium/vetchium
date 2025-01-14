package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrInvalidBase64 = errors.New("invalid base64 encoding")
	ErrNotPDF        = errors.New("file is not a PDF")
	ErrPDFTooLarge   = errors.New("PDF file too large")
	ErrMalformedPDF  = errors.New("malformed PDF file")
)

const (
	maxPDFSize = 10 * 1024 * 1024 // 10MB
	// PDF file header magic number
	pdfHeader = "%PDF-"
)

// ValidateAndSanitizePDF checks if the given base64 string is a valid PDF and performs basic security checks
// Returns the decoded PDF bytes if valid, or an error if invalid
//
// TODO: Additional security measures to consider:
// 1. Use a PDF parsing library (like pdfcpu) to:
//   - Validate complete PDF structure
//   - Check for and disable JavaScript content
//   - Remove embedded files/attachments
//   - Sanitize potentially malicious content
//
// 2. Implement virus scanning:
//   - Integrate with ClamAV or similar antivirus
//   - Scan for known PDF exploits
//   - Check for suspicious patterns
//
// 3. Advanced PDF validation:
//   - Verify PDF version compatibility
//   - Check for encrypted content
//   - Validate digital signatures if present
//   - Scan for malicious URL patterns
//
// 4. Content restrictions:
//   - Limit number of pages
//   - Restrict embedded fonts
//   - Control image resolution/quality
//   - Remove metadata if needed
func ValidateAndSanitizePDF(base64PDF string) ([]byte, error) {
	// Decode base64
	pdfBytes, err := base64.StdEncoding.DecodeString(base64PDF)
	if err != nil {
		return nil, ErrInvalidBase64
	}

	// Check file size
	if len(pdfBytes) > maxPDFSize {
		return nil, ErrPDFTooLarge
	}

	// Check PDF header magic number
	if !bytes.HasPrefix(pdfBytes, []byte(pdfHeader)) {
		return nil, ErrNotPDF
	}

	// Basic PDF structure validation
	err = validatePDFStructure(pdfBytes)
	if err != nil {
		return nil, err
	}

	return pdfBytes, nil
}

// validatePDFStructure performs basic structural validation of a PDF file
// This includes checking for:
// 1. Valid header
// 2. Presence of EOF marker
// 3. Basic xref table structure
func validatePDFStructure(pdfBytes []byte) error {
	r := bytes.NewReader(pdfBytes)

	// Check for EOF marker
	_, err := r.Seek(-6, io.SeekEnd) // "%%EOF" + possible newline
	if err != nil {
		return ErrMalformedPDF
	}

	eofBuf := make([]byte, 6)
	_, err = r.Read(eofBuf)
	if err != nil {
		return ErrMalformedPDF
	}

	if !bytes.Contains(eofBuf, []byte("%%EOF")) {
		return ErrMalformedPDF
	}

	// Check for xref table (a basic PDF structure element)
	if !bytes.Contains(pdfBytes, []byte("xref")) {
		return ErrMalformedPDF
	}

	return nil
}
