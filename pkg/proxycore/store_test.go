package proxycore

import "testing"

func TestStoreEncryptDecrypt(t *testing.T) {
	store := NewStore(t.TempDir())

	plain := "my-secret-password"
	enc, err := store.encrypt(plain)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if enc == plain {
		t.Error("encrypted should differ from plaintext")
	}
	if len(enc) < 5 || enc[:4] != "enc:" {
		t.Errorf("should start with 'enc:', got %q", enc)
	}

	dec, err := store.decrypt(enc)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if dec != plain {
		t.Errorf("decrypted = %q, want %q", dec, plain)
	}
}

func TestStoreEncryptEmpty(t *testing.T) {
	store := NewStore(t.TempDir())

	enc, err := store.encrypt("")
	if err != nil {
		t.Fatal(err)
	}
	if enc != "" {
		t.Errorf("encrypt empty = %q, want empty", enc)
	}

	dec, err := store.decrypt("")
	if err != nil {
		t.Fatal(err)
	}
	if dec != "" {
		t.Errorf("decrypt empty = %q, want empty", dec)
	}
}

func TestStoreDecryptPlaintext(t *testing.T) {
	store := NewStore(t.TempDir())

	dec, err := store.decrypt("plaintext-value")
	if err != nil {
		t.Fatal(err)
	}
	if dec != "plaintext-value" {
		t.Errorf("decrypt plaintext = %q, want original", dec)
	}
}

func TestStoreSaveLoad(t *testing.T) {
	store := NewStore(t.TempDir())

	entries := []*ProxyEntry{
		{
			ID: "t1", Host: "p1.example.com", Port: "8080",
			Username: "user1", Password: "pass1",
			Protocol: "http", Status: StatusConnected, Active: true,
		},
		{
			ID: "t2", Host: "p2.example.com", Port: "3128",
			Protocol: "http", Status: StatusOffline,
		},
	}

	if err := store.Save(entries); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("loaded %d, want 2", len(loaded))
	}

	if loaded[0].Username != "user1" || loaded[0].Password != "pass1" {
		t.Errorf("credentials lost: user=%q pass=%q", loaded[0].Username, loaded[0].Password)
	}
	if loaded[0].Status != StatusOffline {
		t.Errorf("status should reset to offline on load, got %q", loaded[0].Status)
	}
	if loaded[0].Active {
		t.Error("active should reset to false on load")
	}
}

func TestStoreLoadNonExistent(t *testing.T) {
	store := NewStore(t.TempDir())

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load non-existent should not error: %v", err)
	}
	if entries != nil {
		t.Errorf("should return nil for non-existent file, got %d entries", len(entries))
	}
}

func TestStoreKeyPersistence(t *testing.T) {
	dir := t.TempDir()
	s1 := NewStore(dir)
	s2 := NewStore(dir)

	enc, _ := s1.encrypt("secret")
	dec, err := s2.decrypt(enc)
	if err != nil {
		t.Fatalf("cross-store decrypt: %v", err)
	}
	if dec != "secret" {
		t.Errorf("cross-store decrypt = %q, want 'secret'", dec)
	}
}
