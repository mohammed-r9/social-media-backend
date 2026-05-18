package stateful

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateOpaqueToken(length int) string {
	randomByte := make([]byte, length)
	rand.Read(randomByte)
	return hex.EncodeToString(randomByte)
}
