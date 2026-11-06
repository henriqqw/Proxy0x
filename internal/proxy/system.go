package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"proxy_ui/internal/fsutil"

	"golang.org/x/sys/windows/registry"
)

var wininet = syscall.NewLazyDLL("wininet.dll")
var internetSetOption = wininet.NewProc("InternetSetOptionW")

const (
	internetSettingsKey = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`
	bypassList          = "localhost;127.*;10.*;192.168.*;<local>"
)

// Settings are the WinINET proxy values under HKCU\...\Internet Settings.
// An empty Server or Override means the value is absent.
type Settings struct {
	Enable   bool   `json:"enable"`
	Server   string `json:"server"`
	Override string `json:"override"`
}

type Registry interface {
	Get() (Settings, error)
	Set(Settings) error
}

// SystemProxy points the Windows proxy at the local server.
type SystemProxy struct {
	reg        Registry
	refresh    func()
	backupPath string
}

// backup is what the user had before Proxy0x changed the system proxy,
// plus the address Proxy0x set (to recognize our own leftover after a crash).
type backup struct {
	Previous Settings `json:"previous"`
	Ours     string   `json:"ours"`
}

func NewSystemProxy(backupPath string) *SystemProxy {
	return &SystemProxy{reg: winRegistry{}, refresh: refreshProxy, backupPath: backupPath}
}

func (s *SystemProxy) Set(port int) error {
	ours := fmt.Sprintf("127.0.0.1:%d", port)
	b, hasBackup, err := s.readBackup()
	if err != nil {
		return err
	}
	if !hasBackup {
		cur, err := s.reg.Get()
		if err != nil {
			return err
		}
		b.Previous = cur
	}
	b.Ours = ours
	if err := s.writeBackup(b); err != nil {
		return fmt.Errorf("save system proxy backup: %w", err)
	}
	err = s.reg.Set(Settings{Enable: true, Server: ours, Override: bypassList})
	s.refresh()
	return err
}

// Clear restores the user's previous settings.
func (s *SystemProxy) Clear() error {
	b, hasBackup, err := s.readBackup()
	if err != nil {
		return err
	}
	if !hasBackup {
		return s.disableIfLegacyLeftover()
	}
	if err := s.reg.Set(b.Previous); err != nil {
		return err
	}
	s.refresh()
	return os.Remove(s.backupPath)
}

// Recover runs at startup: it undoes only a proxy that Proxy0x itself left
// behind. Settings the user changed since then are never touched.
func (s *SystemProxy) Recover() error {
	b, hasBackup, err := s.readBackup()
	if err != nil {
		return err
	}
	if !hasBackup {
		return s.disableIfLegacyLeftover()
	}
	cur, err := s.reg.Get()
	if err != nil {
		return err
	}
	if cur.Enable && cur.Server == b.Ours {
		return s.Clear()
	}
	return os.Remove(s.backupPath)
}

// disableIfLegacyLeftover handles proxies left by versions without a backup
// file: 127.0.0.1 with Proxy0x's exact bypass list.
func (s *SystemProxy) disableIfLegacyLeftover() error {
	cur, err := s.reg.Get()
	if err != nil {
		return err
	}
	if !cur.Enable || !strings.HasPrefix(cur.Server, "127.0.0.1:") || cur.Override != bypassList {
		return nil
	}
	cur.Enable = false
	err = s.reg.Set(cur)
	s.refresh()
	return err
}

func (s *SystemProxy) readBackup() (backup, bool, error) {
	data, err := os.ReadFile(s.backupPath)
	if os.IsNotExist(err) {
		return backup{}, false, nil
	}
	if err != nil {
		return backup{}, false, err
	}
	var b backup
	if err := json.Unmarshal(data, &b); err != nil {
		return backup{}, false, fmt.Errorf("parse %s: %w", s.backupPath, err)
	}
	return b, true, nil
}

func (s *SystemProxy) writeBackup(b backup) error {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	return fsutil.WriteFileAtomic(s.backupPath, data, 0o600)
}

type winRegistry struct{}

func (winRegistry) Get() (Settings, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, internetSettingsKey, registry.QUERY_VALUE)
	if err != nil {
		return Settings{}, err
	}
	defer k.Close()
	var s Settings
	if v, _, err := k.GetIntegerValue("ProxyEnable"); err == nil {
		s.Enable = v != 0
	} else if !errors.Is(err, registry.ErrNotExist) {
		return Settings{}, err
	}
	if s.Server, err = getString(k, "ProxyServer"); err != nil {
		return Settings{}, err
	}
	if s.Override, err = getString(k, "ProxyOverride"); err != nil {
		return Settings{}, err
	}
	return s, nil
}

func (winRegistry) Set(s Settings) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, internetSettingsKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if err := setString(k, "ProxyServer", s.Server); err != nil {
		return err
	}
	if err := setString(k, "ProxyOverride", s.Override); err != nil {
		return err
	}
	var enable uint32
	if s.Enable {
		enable = 1
	}
	return k.SetDWordValue("ProxyEnable", enable)
}

func getString(k registry.Key, name string) (string, error) {
	v, _, err := k.GetStringValue(name)
	if errors.Is(err, registry.ErrNotExist) {
		return "", nil
	}
	return v, err
}

func setString(k registry.Key, name, value string) error {
	if value == "" {
		if err := k.DeleteValue(name); err != nil && !errors.Is(err, registry.ErrNotExist) {
			return err
		}
		return nil
	}
	return k.SetStringValue(name, value)
}

func refreshProxy() {
	internetSetOption.Call(0, 39, uintptr(unsafe.Pointer(nil)), 0)
	internetSetOption.Call(0, 37, uintptr(unsafe.Pointer(nil)), 0)
}
