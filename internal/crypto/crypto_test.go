package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCorruptKeyIsAnErrorNotReplaced(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, ".key")
	os.WriteFile(keyPath, []byte("garbage-not-a-key"), 0o600)

	if _, err := NewStore(dir).Encrypt("x"); err == nil {
		t.Fatal("a corrupt key must be an error")
	}
	if got, _ := os.ReadFile(keyPath); string(got) != "garbage-not-a-key" {
		t.Fatal("corrupt key file was replaced: existing credentials would be lost for good")
	}
}

func TestDPAPIRoundTrip(t *testing.T) {
	enc, err := dpapiEncrypt([]byte("secret-key-material"))
	if err != nil {
		t.Fatal(err)
	}
	dec, err := dpapiDecrypt(enc)
	if err != nil {
		t.Fatal(err)
	}
	if string(dec) != "secret-key-material" {
		t.Fatalf("round trip = %q", dec)
	}
}

func TestStoreEncryptDecryptWithDPAPIKey(t *testing.T) {
	s := NewStore(t.TempDir())
	enc, err := s.Encrypt("p4ss")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := NewStore(s.dataDir).Decrypt(enc)
	if err != nil || dec != "p4ss" {
		t.Fatalf("Decrypt = %q, %v", dec, err)
	}
}
