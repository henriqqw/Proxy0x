package proxycore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	dir string
	mu  sync.Mutex
}

func NewStore(dataDir string) *Store {
	return &Store{dir: dataDir}
}

func (s *Store) DataDir() string {
	return s.dir
}

func (s *Store) proxiesFile() string {
	return filepath.Join(s.dir, "proxies.json")
}

func (s *Store) keyFile() string {
	return filepath.Join(s.dir, ".key")
}

func (s *Store) getOrCreateKey() ([]byte, error) {
	os.MkdirAll(s.dir, 0o700)
	path := s.keyFile()
	key, err := os.ReadFile(path)
	if err == nil && len(key) == 32 {
		return key, nil
	}
	key = make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, key, 0o600); err != nil {
		return nil, err
	}
	return key, nil
}

func (s *Store) encrypt(plaintext string) (string, error) {
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

func (s *Store) decrypt(encoded string) (string, error) {
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

func (s *Store) Load() ([]*ProxyEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.proxiesFile())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var entries []*ProxyEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	for _, p := range entries {
		if dec, err := s.decrypt(p.Username); err == nil {
			p.Username = dec
		}
		if dec, err := s.decrypt(p.Password); err == nil {
			p.Password = dec
		}
		p.Status = StatusOffline
		p.Active = false
		p.CurrentIP = ""
		p.Latency = 0
	}
	return entries, nil
}

func (s *Store) Save(proxies []*ProxyEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	os.MkdirAll(s.dir, 0o700)

	entries := make([]*ProxyEntry, len(proxies))
	for i, p := range proxies {
		clone := *p
		if enc, err := s.encrypt(clone.Username); err == nil {
			clone.Username = enc
		}
		if enc, err := s.encrypt(clone.Password); err == nil {
			clone.Password = enc
		}
		clone.Status = StatusOffline
		clone.Active = false
		clone.CurrentIP = ""
		clone.Latency = 0
		entries[i] = &clone
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.proxiesFile(), data, 0o600)
}
