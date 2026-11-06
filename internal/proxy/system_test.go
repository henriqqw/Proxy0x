package proxy

import (
	"os"
	"path/filepath"
	"testing"
)

type fakeRegistry struct{ s Settings }

func (f *fakeRegistry) Get() (Settings, error) { return f.s, nil }
func (f *fakeRegistry) Set(s Settings) error   { f.s = s; return nil }

func newTestSystemProxy(t *testing.T, initial Settings) (*SystemProxy, *fakeRegistry, string) {
	t.Helper()
	reg := &fakeRegistry{s: initial}
	backup := filepath.Join(t.TempDir(), "sysproxy-backup.json")
	return &SystemProxy{reg: reg, refresh: func() {}, backupPath: backup}, reg, backup
}

var corporate = Settings{Enable: true, Server: "corp.example:8080", Override: "*.corp;<local>"}

func TestConnectDisconnectRestoresUserProxy(t *testing.T) {
	sp, reg, backup := newTestSystemProxy(t, corporate)

	if err := sp.Set(5000); err != nil {
		t.Fatal(err)
	}
	if reg.s.Server != "127.0.0.1:5000" || !reg.s.Enable {
		t.Fatalf("system proxy not pointed at the local server: %+v", reg.s)
	}
	if err := sp.Clear(); err != nil {
		t.Fatal(err)
	}
	if reg.s != corporate {
		t.Fatalf("user proxy not restored: got %+v, want %+v", reg.s, corporate)
	}
	if _, err := os.Stat(backup); !os.IsNotExist(err) {
		t.Fatal("backup should be deleted after restore")
	}
}

func TestPortChangeKeepsOriginalBackup(t *testing.T) {
	sp, reg, _ := newTestSystemProxy(t, corporate)
	sp.Set(5000)
	sp.Set(6000) // reconnect / rotation on a new port

	sp.Clear()
	if reg.s != corporate {
		t.Fatalf("restored %+v, want the user's original %+v", reg.s, corporate)
	}
}

func TestRecoverRestoresAfterCrash(t *testing.T) {
	sp, reg, _ := newTestSystemProxy(t, corporate)
	sp.Set(5000)
	// app crashed here; next start:
	sp2 := &SystemProxy{reg: reg, refresh: func() {}, backupPath: sp.backupPath}

	if err := sp2.Recover(); err != nil {
		t.Fatal(err)
	}
	if reg.s != corporate {
		t.Fatalf("recovered %+v, want %+v", reg.s, corporate)
	}
}

func TestRecoverLeavesOtherLocalToolsAlone(t *testing.T) {
	fiddler := Settings{Enable: true, Server: "127.0.0.1:8888", Override: "<-loopback>"}
	sp, reg, _ := newTestSystemProxy(t, fiddler)

	sp.Recover()
	if reg.s != fiddler {
		t.Fatalf("another tool's proxy was changed: %+v", reg.s)
	}
}

func TestRecoverDiscardsBackupWhenUserChangedSettings(t *testing.T) {
	sp, reg, backup := newTestSystemProxy(t, corporate)
	sp.Set(5000)
	userNow := Settings{Enable: true, Server: "other.example:3128"}
	reg.s = userNow // user reconfigured while the app was dead

	sp.Recover()
	if reg.s != userNow {
		t.Fatalf("user's newer settings were overwritten: %+v", reg.s)
	}
	if _, err := os.Stat(backup); !os.IsNotExist(err) {
		t.Fatal("stale backup should be discarded")
	}
}

func TestRecoverDisablesLegacyLeftover(t *testing.T) {
	legacy := Settings{Enable: true, Server: "127.0.0.1:51234", Override: bypassList}
	sp, reg, _ := newTestSystemProxy(t, legacy)

	sp.Recover()
	if reg.s.Enable {
		t.Fatalf("leftover from an older version still enabled: %+v", reg.s)
	}
}
