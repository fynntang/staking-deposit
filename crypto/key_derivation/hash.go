package key_derivation

import (
	"crypto/sha256"

	"golang.org/x/crypto/sha3"
)

// SHA3256 creates an SHA3-256 hash of the supplied data
func SHA3256(data ...[]byte) []byte {
	hash := sha3.New256()
	for _, d := range data {
		_, _ = hash.Write(d)
	}
	return hash.Sum(nil)
}

// SHA256 creates an SHA-256 hash of the supplied data
func SHA256(data ...[]byte) []byte {
	hash := sha256.New()
	for _, d := range data {
		_, _ = hash.Write(d)
	}
	return hash.Sum(nil)
}

// Keccak256 creates a Keccak256 hash of the supplied data
func Keccak256(data ...[]byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	for _, d := range data {
		_, _ = hash.Write(d)
	}
	return hash.Sum(nil)
}
