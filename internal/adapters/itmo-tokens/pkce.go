package itmotokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// generateCodeVerifier creates a random code verifier for PKCE.
func generateCodeVerifier() (string, error) {
	b := make([]byte, 40)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// getCodeChallenge returns a S256 code challenge for PKCE.
func getCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
