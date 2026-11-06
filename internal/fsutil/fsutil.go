package fsutil

import (
	"os"
	"path/filepath"
)

// WriteFileAtomic writes data to path via a temp file and rename, so a crash
// never leaves a truncated file. The previous content is kept as path+".bak".
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	return WriteFileAtomicWithBackup(path, path+".bak", data, perm)
}

// WriteFileAtomicWithBackup is WriteFileAtomic with an explicit backup path.
func WriteFileAtomicWithBackup(path, backup string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	if old, err := os.ReadFile(path); err == nil {
		if err := os.WriteFile(backup, old, perm); err != nil {
			return err
		}
	}

	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}
