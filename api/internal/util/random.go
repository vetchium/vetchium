package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/rs/xid"
)

// TODO: Usage of this should be deprecated and replaced with RandomUniqueID
func RandomString(numBytes int) string {
	buff := make([]byte, numBytes)
	rand.Read(buff)

	// The lower case conversion is added to guard against any regression
	// with hex.EncodeToString returning uppercase.
	return strings.ToLower(hex.EncodeToString(buff))
}

func RandNumString(numDigits int) (string, error) {
	// Generate a 6-digit random number
	min := int64(100000) // Smallest 6-digit number
	max := int64(999999) // Largest 6-digit number
	rangeSize := max - min + 1

	num, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		return "", err
	}

	// Adjust to fall within the 6-digit range
	randomNumber := num.Int64() + min

	// Convert to string
	return fmt.Sprintf("%d", randomNumber), nil
}

// RandomUniqueID generates a random unique ID of the form
// <xid>.<random-string-of-numBytes> The xid is a 12 character string that is
// not cryptographically secure but can help with sorting and uniqueness.
func RandomUniqueID(numBytes int) string {
	buff := make([]byte, numBytes)
	rand.Read(buff)

	// The lower case conversion is added to guard against any regression
	// with hex.EncodeToString returning uppercase.
	return strings.ToLower(xid.New().String() + hex.EncodeToString(buff))
}
