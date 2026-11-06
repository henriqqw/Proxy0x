package killswitch

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"

	"proxy_ui/internal/fsutil"
)

const (
	ruleAllowApp   = "Proxy0x-KillSwitch-AllowApp"
	ruleAllowLocal = "Proxy0x-KillSwitch-AllowLocal"
	// Go resolves names through the DNS Client service (svchost), not from
	// Proxy0x.exe, so upstream hostnames need their own exception.
	ruleAllowDNSUDP = "Proxy0x-KillSwitch-AllowDNS-UDP"
	ruleAllowDNSTCP = "Proxy0x-KillSwitch-AllowDNS-TCP"
	// Created by versions up to 1.1.0. Windows Firewall applies block rules
	// before allow rules, so it also cut off Proxy0x itself.
	ruleLegacyBlock = "Proxy0x-KillSwitch-Block"

	createNoWindow = 0x08000000
)

var profiles = []string{"domain", "private", "public"}

var (
	policyInOutput = regexp.MustCompile(`(?i)\b(blockinboundalways|blockinbound|allowinbound),(blockoutbound|allowoutbound)\b`)
	policyExact    = regexp.MustCompile(`^(blockinboundalways|blockinbound|allowinbound),(blockoutbound|allowoutbound)$`)
)

// Firewall blocks outbound traffic for every program except Proxy0x by setting
// the outbound default of each firewall profile to block and allowing only
// Proxy0x and loopback. A block rule cannot do this: block rules win over
// allow rules. The policy in place before Enable is saved to backupPath and
// restored by Disable.
type Firewall struct {
	backupPath string
	exePath    func() (string, error)
	systemDir  func() (string, error)
	readPolicy func(profile string) (string, error)
	ruleExists func(name string) bool
	elevate    func(command string) error
}

func New(backupPath string) *Firewall {
	return &Firewall{
		backupPath: backupPath,
		exePath:    os.Executable,
		systemDir:  windows.GetSystemDirectory,
		readPolicy: netshPolicy,
		ruleExists: netshRuleExists,
		elevate:    RunElevated,
	}
}

// policies maps a firewall profile to its "inbound,outbound" policy.
type policies map[string]string

func (f *Firewall) Active() bool {
	return f.ruleExists(ruleAllowApp)
}

func (f *Firewall) Enable() error {
	exe, err := f.exePath()
	if err != nil {
		return fmt.Errorf("get exe path: %w", err)
	}
	sysDir, err := f.systemDir()
	if err != nil {
		return fmt.Errorf("locate system directory: %w", err)
	}
	svchost := filepath.Join(sysDir, "svchost.exe")
	for _, p := range []string{exe, svchost} {
		if strings.ContainsAny(p, `"%`) {
			return fmt.Errorf("unsupported characters in program path %q", p)
		}
	}

	saved, hasBackup, err := f.readBackup()
	if err != nil {
		return err
	}
	// An existing backup holds the user's policy from an earlier Enable that
	// was never undone; re-reading now would capture our own blockoutbound.
	if !hasBackup {
		if saved, err = f.currentPolicies(); err != nil {
			return err
		}
		if err := f.writeBackup(saved); err != nil {
			return fmt.Errorf("save firewall policy backup: %w", err)
		}
	}

	if err := f.elevate(enableCommand(exe, svchost, saved)); err != nil {
		if !hasBackup && errors.Is(err, ErrCancelled) {
			os.Remove(f.backupPath)
		}
		return err
	}
	return nil
}

func (f *Firewall) Disable() error {
	restore, hasBackup, err := f.readBackup()
	if err != nil {
		return err
	}
	if !hasBackup {
		if !f.ruleExists(ruleAllowApp) && !f.ruleExists(ruleAllowLocal) && !f.ruleExists(ruleLegacyBlock) {
			return nil
		}
		restore = policies{}
		// Without a backup only undo what Proxy0x most likely did: outbound
		// blocked while its allow rule exists. A user's own policy is left alone.
		if f.ruleExists(ruleAllowApp) {
			current, err := f.currentPolicies()
			if err != nil {
				return err
			}
			for p, pol := range current {
				in, out, _ := strings.Cut(pol, ",")
				if out == "blockoutbound" {
					restore[p] = in + ",allowoutbound"
				}
			}
		}
	}

	if err := f.elevate(disableCommand(restore)); err != nil {
		return err
	}
	if hasBackup {
		if err := os.Remove(f.backupPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func (f *Firewall) currentPolicies() (policies, error) {
	out := policies{}
	for _, p := range profiles {
		pol, err := f.readPolicy(p)
		if err != nil {
			return nil, fmt.Errorf("read %s firewall policy: %w", p, err)
		}
		if !policyExact.MatchString(pol) {
			return nil, fmt.Errorf("unexpected %s firewall policy %q", p, pol)
		}
		out[p] = pol
	}
	return out, nil
}

func (f *Firewall) readBackup() (policies, bool, error) {
	data, err := os.ReadFile(f.backupPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var b policies
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, false, fmt.Errorf("parse %s: %w", f.backupPath, err)
	}
	// The file is user-writable and its values end up in an elevated command.
	for _, p := range profiles {
		if !policyExact.MatchString(b[p]) {
			return nil, false, fmt.Errorf("invalid %s policy %q in %s", p, b[p], f.backupPath)
		}
	}
	return b, true, nil
}

func (f *Firewall) writeBackup(b policies) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return fsutil.WriteFileAtomic(f.backupPath, data, 0o600)
}

// enableCommand adds the allow rules first and blocks outbound last, so a
// failed rule never leaves the machine blocked without an exception for Proxy0x.
func enableCommand(exe, svchost string, saved policies) string {
	steps := []string{
		fmt.Sprintf(`netsh advfirewall firewall add rule name="%s" dir=out action=allow program="%s" profile=any enable=yes`, ruleAllowApp, exe),
		fmt.Sprintf(`netsh advfirewall firewall add rule name="%s" dir=out action=allow remoteip=127.0.0.1,::1 profile=any enable=yes`, ruleAllowLocal),
		fmt.Sprintf(`netsh advfirewall firewall add rule name="%s" dir=out action=allow program="%s" service=dnscache protocol=udp remoteport=53 profile=any enable=yes`, ruleAllowDNSUDP, svchost),
		fmt.Sprintf(`netsh advfirewall firewall add rule name="%s" dir=out action=allow program="%s" service=dnscache protocol=tcp remoteport=53 profile=any enable=yes`, ruleAllowDNSTCP, svchost),
	}
	for _, p := range profiles {
		in, _, _ := strings.Cut(saved[p], ",")
		steps = append(steps, setPolicy(p, in+",blockoutbound"))
	}
	return deleteRules() + " & " + strings.Join(steps, " && ")
}

// disableCommand restores the policy first and removes the rules last, so
// outbound is never blocked once the allow rules are gone.
func disableCommand(restore policies) string {
	var steps []string
	for _, p := range profiles {
		if pol, ok := restore[p]; ok {
			steps = append(steps, setPolicy(p, pol))
		}
	}
	steps = append(steps, "("+deleteRules()+" & exit 0)")
	return strings.Join(steps, " && ")
}

func deleteRules() string {
	var cmds []string
	for _, name := range []string{ruleLegacyBlock, ruleAllowApp, ruleAllowLocal, ruleAllowDNSUDP, ruleAllowDNSTCP} {
		cmds = append(cmds, fmt.Sprintf(`netsh advfirewall firewall delete rule name="%s" >nul 2>&1`, name))
	}
	return strings.Join(cmds, " & ")
}

func setPolicy(profile, policy string) string {
	return fmt.Sprintf("netsh advfirewall set %sprofile firewallpolicy %s", profile, policy)
}

func netsh(args ...string) ([]byte, error) {
	cmd := exec.Command("netsh", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNoWindow}
	return cmd.CombinedOutput()
}

func netshPolicy(profile string) (string, error) {
	out, err := netsh("advfirewall", "show", profile+"profile", "firewallpolicy")
	if err != nil {
		return "", fmt.Errorf("netsh: %w", err)
	}
	m := policyInOutput.FindString(string(out))
	if m == "" {
		return "", fmt.Errorf("no firewall policy in netsh output")
	}
	return strings.ToLower(m), nil
}

// netshRuleExists relies on netsh exiting with 1 when no rule matches.
func netshRuleExists(name string) bool {
	_, err := netsh("advfirewall", "firewall", "show", "rule", "name="+name)
	return err == nil
}
