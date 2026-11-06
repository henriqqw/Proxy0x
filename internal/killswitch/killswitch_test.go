package killswitch

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeFirewall struct {
	policy   map[string]string
	rules    map[string]bool
	commands []string
	runErr   error
}

func newFake(t *testing.T) (*Firewall, *fakeFirewall) {
	t.Helper()
	fw := &fakeFirewall{
		policy: map[string]string{
			"domain":  "blockinbound,allowoutbound",
			"private": "allowinbound,allowoutbound",
			"public":  "blockinboundalways,allowoutbound",
		},
		rules: map[string]bool{},
	}
	f := &Firewall{
		backupPath: filepath.Join(t.TempDir(), "killswitch-backup.json"),
		exePath:    func() (string, error) { return `C:\Program Files\Proxy0x\Proxy0x.exe`, nil },
		systemDir:  func() (string, error) { return `C:\Windows\system32`, nil },
		readPolicy: func(p string) (string, error) { return fw.policy[p], nil },
		ruleExists: func(name string) bool { return fw.rules[name] },
		elevate: func(cmd string) error {
			fw.commands = append(fw.commands, cmd)
			return fw.runErr
		},
	}
	return f, fw
}

func readBackupFile(t *testing.T, path string) policies {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	var b policies
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("parse backup: %v", err)
	}
	return b
}

func TestEnableAllowsProxy0xBeforeBlockingOutbound(t *testing.T) {
	f, fw := newFake(t)
	if err := f.Enable(); err != nil {
		t.Fatal(err)
	}
	if len(fw.commands) != 1 {
		t.Fatalf("want 1 elevated command, got %d", len(fw.commands))
	}
	cmd := fw.commands[0]

	if strings.Contains(cmd, "action=block") {
		t.Errorf("must not add a block rule (it overrides the allow rule): %s", cmd)
	}
	allowApp := strings.Index(cmd, `name="Proxy0x-KillSwitch-AllowApp" dir=out action=allow program="C:\Program Files\Proxy0x\Proxy0x.exe"`)
	allowLocal := strings.Index(cmd, `name="Proxy0x-KillSwitch-AllowLocal" dir=out action=allow remoteip=127.0.0.1,::1`)
	allowDNS := strings.Index(cmd, `name="Proxy0x-KillSwitch-AllowDNS-UDP" dir=out action=allow program="C:\Windows\system32\svchost.exe" service=dnscache protocol=udp remoteport=53`)
	firstBlock := strings.Index(cmd, "blockoutbound")
	if allowApp < 0 || allowLocal < 0 || allowDNS < 0 || firstBlock < 0 {
		t.Fatalf("missing allow rules or policy change: %s", cmd)
	}
	if !strings.Contains(cmd, `name="Proxy0x-KillSwitch-AllowDNS-TCP" dir=out action=allow program="C:\Windows\system32\svchost.exe" service=dnscache protocol=tcp remoteport=53`) {
		t.Errorf("missing TCP DNS rule: %s", cmd)
	}
	if firstBlock < allowApp || firstBlock < allowLocal || firstBlock < allowDNS {
		t.Errorf("outbound must be blocked only after the allow rules: %s", cmd)
	}
	for _, want := range []string{
		"set domainprofile firewallpolicy blockinbound,blockoutbound",
		"set privateprofile firewallpolicy allowinbound,blockoutbound",
		"set publicprofile firewallpolicy blockinboundalways,blockoutbound",
		`delete rule name="Proxy0x-KillSwitch-Block"`,
	} {
		if !strings.Contains(cmd, want) {
			t.Errorf("command lacks %q: %s", want, cmd)
		}
	}
	// Adds and policy changes are chained with && so a failure stops the rest.
	if strings.Contains(cmd[allowApp:], " & ") {
		t.Errorf("steps after the deletes must be chained with &&: %s", cmd)
	}
}

func TestEnableSavesOriginalPolicyOnce(t *testing.T) {
	f, fw := newFake(t)
	if err := f.Enable(); err != nil {
		t.Fatal(err)
	}
	want := policies{
		"domain":  "blockinbound,allowoutbound",
		"private": "allowinbound,allowoutbound",
		"public":  "blockinboundalways,allowoutbound",
	}
	if got := readBackupFile(t, f.backupPath); len(got) != 3 || got["domain"] != want["domain"] || got["private"] != want["private"] || got["public"] != want["public"] {
		t.Fatalf("backup = %v, want %v", got, want)
	}

	// The firewall now blocks outbound; enabling again must keep the user's policy.
	for p := range fw.policy {
		in, _, _ := strings.Cut(fw.policy[p], ",")
		fw.policy[p] = in + ",blockoutbound"
	}
	if err := f.Enable(); err != nil {
		t.Fatal(err)
	}
	if got := readBackupFile(t, f.backupPath); got["domain"] != want["domain"] {
		t.Fatalf("second Enable overwrote the backup: %v", got)
	}
}

func TestEnableCancelledLeavesNoBackup(t *testing.T) {
	f, fw := newFake(t)
	fw.runErr = ErrCancelled
	if err := f.Enable(); !errors.Is(err, ErrCancelled) {
		t.Fatalf("err = %v, want ErrCancelled", err)
	}
	if _, err := os.Stat(f.backupPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("backup should be removed after a cancelled UAC prompt, stat err = %v", err)
	}
}

func TestEnableFailureAfterElevationKeepsBackup(t *testing.T) {
	f, fw := newFake(t)
	fw.runErr = errors.New("elevated command failed (exit code 1)")
	if err := f.Enable(); err == nil {
		t.Fatal("want error")
	}
	if _, err := os.Stat(f.backupPath); err != nil {
		t.Fatalf("backup must survive a partial run so Disable can restore it: %v", err)
	}
}

func TestEnableRejectsUnsafeExePath(t *testing.T) {
	f, fw := newFake(t)
	f.exePath = func() (string, error) { return `C:\%TEMP%\Proxy0x.exe`, nil }
	if err := f.Enable(); err == nil {
		t.Fatal("want error for a path cmd.exe would expand")
	}
	if len(fw.commands) != 0 {
		t.Fatalf("nothing may run, got %v", fw.commands)
	}
}

func TestDisableRestoresBackupThenDeletesRules(t *testing.T) {
	f, fw := newFake(t)
	if err := f.Enable(); err != nil {
		t.Fatal(err)
	}
	fw.commands = nil

	if err := f.Disable(); err != nil {
		t.Fatal(err)
	}
	cmd := fw.commands[0]
	restore := strings.Index(cmd, "set domainprofile firewallpolicy blockinbound,allowoutbound")
	del := strings.Index(cmd, `delete rule name="Proxy0x-KillSwitch-AllowApp"`)
	if !strings.Contains(cmd, `delete rule name="Proxy0x-KillSwitch-AllowDNS-UDP"`) || !strings.Contains(cmd, `delete rule name="Proxy0x-KillSwitch-AllowDNS-TCP"`) {
		t.Errorf("DNS rules must be removed too: %s", cmd)
	}
	if restore < 0 || del < 0 || del < restore {
		t.Fatalf("policy must be restored before the allow rules are deleted: %s", cmd)
	}
	for _, want := range []string{
		"set privateprofile firewallpolicy allowinbound,allowoutbound",
		"set publicprofile firewallpolicy blockinboundalways,allowoutbound",
		"exit 0",
	} {
		if !strings.Contains(cmd, want) {
			t.Errorf("command lacks %q: %s", want, cmd)
		}
	}
	if _, err := os.Stat(f.backupPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("backup should be removed after Disable, stat err = %v", err)
	}
}

func TestDisableFailureKeepsBackup(t *testing.T) {
	f, fw := newFake(t)
	if err := f.Enable(); err != nil {
		t.Fatal(err)
	}
	fw.runErr = ErrCancelled
	if err := f.Disable(); !errors.Is(err, ErrCancelled) {
		t.Fatalf("err = %v, want ErrCancelled", err)
	}
	if _, err := os.Stat(f.backupPath); err != nil {
		t.Fatalf("backup must survive a failed Disable: %v", err)
	}
}

func TestDisableWithoutBackupOrRulesDoesNothing(t *testing.T) {
	f, fw := newFake(t)
	if err := f.Disable(); err != nil {
		t.Fatal(err)
	}
	if len(fw.commands) != 0 {
		t.Fatalf("no UAC prompt when there is nothing to undo, got %v", fw.commands)
	}
}

func TestDisableWithoutBackupLeavesUserPolicyAlone(t *testing.T) {
	f, fw := newFake(t)
	fw.rules[ruleLegacyBlock] = true // only the old block rule is left
	fw.policy["domain"] = "blockinbound,blockoutbound"

	if err := f.Disable(); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(fw.commands[0], "firewallpolicy") {
		t.Fatalf("without our allow rule the policy is the user's own: %s", fw.commands[0])
	}
}

func TestDisableWithoutBackupUnblocksOutbound(t *testing.T) {
	f, fw := newFake(t)
	fw.rules[ruleAllowApp] = true
	fw.policy["domain"] = "blockinbound,blockoutbound"

	if err := f.Disable(); err != nil {
		t.Fatal(err)
	}
	cmd := fw.commands[0]
	if !strings.Contains(cmd, "set domainprofile firewallpolicy blockinbound,allowoutbound") {
		t.Errorf("blocked profile should be reopened: %s", cmd)
	}
	if strings.Contains(cmd, "set privateprofile") {
		t.Errorf("profiles that were not blocked must not change: %s", cmd)
	}
}

func TestTamperedBackupIsRejected(t *testing.T) {
	f, fw := newFake(t)
	bad := `{"domain":"blockinbound,allowoutbound & calc.exe","private":"allowinbound,allowoutbound","public":"blockinbound,allowoutbound"}`
	if err := os.WriteFile(f.backupPath, []byte(bad), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := f.Disable(); err == nil {
		t.Fatal("want error for an invalid policy in the backup")
	}
	if err := f.Enable(); err == nil {
		t.Fatal("want error for an invalid policy in the backup")
	}
	if len(fw.commands) != 0 {
		t.Fatalf("a tampered backup must never reach the elevated command: %v", fw.commands)
	}
}

func TestPolicyParsedFromNetshOutput(t *testing.T) {
	out := "Domain Profile Settings: \r\n------------------\r\nFirewall Policy                       BlockInbound,AllowOutbound\r\nOk.\r\n"
	if got := strings.ToLower(policyInOutput.FindString(out)); got != "blockinbound,allowoutbound" {
		t.Fatalf("parsed %q", got)
	}
}
