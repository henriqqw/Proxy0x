package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

func dpapiEncrypt(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	in := windows.DataBlob{Size: uint32(len(plaintext)), Data: &plaintext[0]}
	var out windows.DataBlob
	if err := windows.CryptProtectData(&in, nil, nil, 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return nil, fmt.Errorf("CryptProtectData failed: %w", err)
	}
	return takeBlob(out), nil
}

func dpapiDecrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	in := windows.DataBlob{Size: uint32(len(ciphertext)), Data: &ciphertext[0]}
	var out windows.DataBlob
	if err := windows.CryptUnprotectData(&in, nil, nil, 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return nil, fmt.Errorf("CryptUnprotectData failed: %w", err)
	}
	return takeBlob(out), nil
}

// takeBlob copies a DPAPI output blob into Go memory and frees it.
func takeBlob(b windows.DataBlob) []byte {
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(b.Data)))
	return append([]byte(nil), unsafe.Slice(b.Data, b.Size)...)
}

type Store struct {
	dataDir string
}

func NewStore(dataDir string) *Store {
	return &Store{dataDir: dataDir}
}

func (s *Store) keyFile() string {
	return filepath.Join(s.dataDir, ".key")
}

func (s *Store) getOrCreateKey() ([]byte, error) {
	os.MkdirAll(s.dataDir, 0o700)
	path := s.keyFile()

	data, err := os.ReadFile(path)
	if err == nil {
		// DPAPI-protected key: "dpapi:" + base64
		if strings.HasPrefix(string(data), "dpapi:") {
			encoded := string(data)[6:]
			encrypted, decErr := base64.StdEncoding.DecodeString(encoded)
			if decErr != nil {
				return nil, fmt.Errorf("decode dpapi key: %w", decErr)
			}
			key, decErr := dpapiDecrypt(encrypted)
			if decErr != nil {
				return nil, fmt.Errorf("dpapi decrypt key (profile changed?): %w — restore from %s.legacy if available", decErr, path)
			}
			if len(key) == 32 {
				return key, nil
			}
			return nil, fmt.Errorf("dpapi key wrong size: %d", len(key))
		}

		// Legacy raw key file — migrate to DPAPI
		if len(data) == 32 {
			key := make([]byte, 32)
			copy(key, data)
			// Keep a backup before overwriting
			os.WriteFile(path+".legacy", data, 0o600)
			if saveErr := s.saveKeyDPAPI(key); saveErr != nil {
				return key, nil
			}
			return key, nil
		}
	}

	// Generate new key, store with DPAPI
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	if err := s.saveKeyDPAPI(key); err != nil {
		return key, os.WriteFile(path, key, 0o600)
	}
	return key, nil
}

func (s *Store) saveKeyDPAPI(key []byte) error {
	encrypted, err := dpapiEncrypt(key)
	if err != nil {
		return err
	}
	encoded := "dpapi:" + base64.StdEncoding.EncodeToString(encrypted)
	return os.WriteFile(s.keyFile(), []byte(encoded), 0o600)
}

func (s *Store) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	key, err := s.getOrCreateKey()
	if err != nil {
		return plaintext, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return plaintext, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return plaintext, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return plaintext, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return "enc:" + base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *Store) Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	if len(encoded) < 5 || encoded[:4] != "enc:" {
		return encoded, nil
	}

	key, err := s.getOrCreateKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(encoded[4:])
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
