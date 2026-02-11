package security

import "testing"

func TestSign(t *testing.T) {
	secret := "test-secret"
	checksumKey := "checksum-key"
	signer := NewSigner(secret, checksumKey)

	data := "hello-world"
	signature := signer.Sign(data)

	if signature == "" {
		t.Fatal("signature should not be empty")
	}
}

func TestVerify(t *testing.T) {
	secret := "test-secret"
	checksumKey := "checksum-key"
	signer := NewSigner(secret, checksumKey)

	data := "hello-world"
	signature := signer.Sign(data)

	if !signer.Verify(data, signature) {
		t.Fatal("verify should return true for valid signature")
	}

	if signer.Verify(data, "invalid-signature") {
		t.Fatal("verify should return false for invalid signature")
	}
}

func TestChecksum(t *testing.T) {
	secret := "test-secret"
	checksumKey := "checksum-key"
	signer := NewSigner(secret, checksumKey)

	data := "sample-data"
	checksum := signer.Checksum(data)

	if checksum == "" {
		t.Fatal("checksum should not be empty")
	}

	// SHA256 hex uppercase => length must be 64
	if len(checksum) != 64 {
		t.Fatalf("checksum length expected 64, got %d", len(checksum))
	}
}

func TestVerifyChecksum(t *testing.T) {
	secret := "test-secret"
	checksumKey := "checksum-key"
	signer := NewSigner(secret, checksumKey)

	data := "sample-data"
	checksum := signer.Checksum(data)

	if !signer.VerifyChecksum(data, checksum) {
		t.Fatal("VerifyChecksum should return true for valid checksum")
	}

	if signer.VerifyChecksum(data, "INVALIDCHECKSUM") {
		t.Fatal("VerifyChecksum should return false for invalid checksum")
	}
}

func TestVerifyChecksum_EmptyInput(t *testing.T) {
	signer := NewSigner("secret", "checksum")

	if signer.VerifyChecksum("", "abc") {
		t.Fatal("VerifyChecksum should return false when data is empty")
	}

	if signer.VerifyChecksum("data", "") {
		t.Fatal("VerifyChecksum should return false when checksum is empty")
	}
}
