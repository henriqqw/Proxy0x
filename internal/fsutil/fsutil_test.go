package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomicKeepsBackupAndLeavesNoTemp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.json")

	if err := WriteFileAtomic(path, []byte("v1"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteFileAtomic(path, []byte("v2"), 0o600); err != nil {
		t.Fatal(err)
	}

	if got, _ := os.ReadFile(path); string(got) != "v2" {
		t.Fatalf("file = %q, want v2", got)
	}
	if got, _ := os.ReadFile(path + ".bak"); string(got) != "v1" {
		t.Fatalf("backup = %q, want v1", got)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatalf("temp file left behind: %v", err)
	}
}

func TestWriteFileAtomicCreatesParentDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "data.json")
	if err := WriteFileAtomic(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path + ".bak"); !os.IsNotExist(err) {
		t.Fatal("no backup expected on first write")
	}
}
