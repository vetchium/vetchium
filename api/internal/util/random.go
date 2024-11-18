package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

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
