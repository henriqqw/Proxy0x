package proxycore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"proxy_ui/internal/crypto"
	"proxy_ui/internal/fsutil"
)

// Store reads the desktop app's proxies.json. Encryption (key file, DPAPI,
// format) is shared with the app through internal/crypto: one implementation,
// so the MCP server can never rewrite the app's key.
type Store struct {
	dir   string
	mu    sync.Mutex
	crypt *crypto.Store
}

func NewStore(dataDir string) *Store {
	return &Store{dir: dataDir, crypt: crypto.NewStore(dataDir)}
}

func (s *Store) DataDir() string {
	return s.dir
}

func (s *Store) proxiesFile() string {
	return filepath.Join(s.dir, "proxies.json")
}

func (s *Store) encrypt(plaintext string) (string, error) { return s.crypt.Encrypt(plaintext) }
func (s *Store) decrypt(encoded string) (string, error)   { return s.crypt.Decrypt(encoded) }

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
		user, errU := s.decrypt(p.Username)
		pass, errP := s.decrypt(p.Password)
		if err := errors.Join(errU, errP); err != nil {
			return nil, fmt.Errorf("decrypt credentials of %s:%s: %w", p.Host, p.Port, err)
		}
		p.Username, p.Password = user, pass
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

	entries := make([]*ProxyEntry, len(proxies))
	for i, p := range proxies {
		clone := *p
		user, errU := s.encrypt(clone.Username)
		pass, errP := s.encrypt(clone.Password)
		if err := errors.Join(errU, errP); err != nil {
			return fmt.Errorf("encrypt credentials: %w", err)
		}
		clone.Username, clone.Password = user, pass
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
	return fsutil.WriteFileAtomic(s.proxiesFile(), data, 0o600)
}
