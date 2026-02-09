package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

type Signer struct {
	secret      []byte
	checksumKey []byte
}

func NewSigner(secret, checksumKey string) *Signer {
	return &Signer{
		secret:      []byte(secret),
		checksumKey: []byte(checksumKey),
	}
}

func (s *Signer) Sign(data string) string {
	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (s *Signer) Verify(data, signature string) bool {
	expected := s.Sign(data)
	return expected == signature
}

func (s *Signer) Checksum(data string) string {
	hash := sha256.Sum256(append([]byte(data), s.checksumKey...))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func (s *Signer) VerifyChecksum(data, checksum string) bool {
	if data == "" || checksum == "" {
		return false
	}

	expected := s.Checksum(data)

	return subtle.ConstantTimeCompare(
		[]byte(expected),
		[]byte(checksum),
	) == 1
}
