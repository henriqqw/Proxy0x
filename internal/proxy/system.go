package proxy

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

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
	reg     Registry
	refresh func()
}

func NewSystemProxy() *SystemProxy {
	return &SystemProxy{reg: winRegistry{}, refresh: refreshProxy}
}

func (s *SystemProxy) Set(port int) error {
	err := s.reg.Set(Settings{Enable: true, Server: fmt.Sprintf("127.0.0.1:%d", port), Override: bypassList})
	s.refresh()
	return err
}

func (s *SystemProxy) Clear() error {
	cur, err := s.reg.Get()
	if err != nil {
		return err
	}
	cur.Enable = false
	err = s.reg.Set(cur)
	s.refresh()
	return err
}

// Recover runs at startup to undo a proxy left behind by a crash.
func (s *SystemProxy) Recover() error {
	return s.Clear()
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
