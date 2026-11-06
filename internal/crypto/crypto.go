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
	"syscall"
	"unsafe"
)

var (
	crypt32            = syscall.NewLazyDLL("crypt32.dll")
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	cryptProtectData   = crypt32.NewProc("CryptProtectData")
	cryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
	localFree          = kernel32.NewProc("LocalFree")
)

type dataBlob struct {
	size uint32
	data uintptr
}

func dpapiEncrypt(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	inBlob := dataBlob{size: uint32(len(plaintext)), data: uintptr(unsafe.Pointer(&plaintext[0]))}
	var outBlob dataBlob

	r, _, err := cryptProtectData.Call(
		uintptr(unsafe.Pointer(&inBlob)),
		0, 0, 0, 0,
		1, // CRYPTPROTECT_UI_FORBIDDEN
		uintptr(unsafe.Pointer(&outBlob)),
	)
	if r == 0 {
		return nil, fmt.Errorf("CryptProtectData failed: %v", err)
	}
	defer localFree.Call(outBlob.data)

	result := make([]byte, outBlob.size)
	copy(result, unsafe.Slice((*byte)(unsafe.Pointer(outBlob.data)), outBlob.size))
	return result, nil
}

func dpapiDecrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	inBlob := dataBlob{size: uint32(len(ciphertext)), data: uintptr(unsafe.Pointer(&ciphertext[0]))}
	var outBlob dataBlob

	r, _, err := cryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&inBlob)),
		0, 0, 0, 0,
		1, // CRYPTPROTECT_UI_FORBIDDEN
		uintptr(unsafe.Pointer(&outBlob)),
	)
	if r == 0 {
		return nil, fmt.Errorf("CryptUnprotectData failed: %v", err)
	}
	defer localFree.Call(outBlob.data)

	result := make([]byte, outBlob.size)
	copy(result, unsafe.Slice((*byte)(unsafe.Pointer(outBlob.data)), outBlob.size))
	return result, nil
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
