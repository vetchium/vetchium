package util

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func RandomString(numBytes int) string {
	buff := make([]byte, numBytes)
	rand.Read(buff)
	return strings.ToLower(hex.EncodeToString(buff))
}
