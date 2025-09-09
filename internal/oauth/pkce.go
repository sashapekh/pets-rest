package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func randBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

func b64url(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func NewPKCE() (codeVerif, challenge string) {
	v := b64url(randBytes(32))
	sum := sha256.Sum256([]byte(v))
	return v, b64url(sum[:])
}

func NewState() string {
	return b64url(randBytes(24))
}
