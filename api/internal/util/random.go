package util

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomString(numBytes int) string {
	buff := make([]byte, numBytes)
	rand.Read(buff)
	return hex.EncodeToString(buff)
}
