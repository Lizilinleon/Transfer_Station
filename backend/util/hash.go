package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashSHA256 produces a deterministic SHA-256 hex digest.
func HashSHA256(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}
