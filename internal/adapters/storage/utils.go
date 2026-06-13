package storage

import (
	"crypto/rand"
	"encoding/hex"
)

func GenereateObjectKey() string {
	bytes := make([]byte, 12)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
