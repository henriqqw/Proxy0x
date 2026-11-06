package crypto

import "testing"

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
