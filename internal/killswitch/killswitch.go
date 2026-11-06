package killswitch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	shell32      = syscall.NewLazyDLL("shell32.dll")
	shellExecute = shell32.NewProc("ShellExecuteW")
)

const (
	fwRuleBlock      = "Proxy0x-KillSwitch-Block"
	fwRuleAllowApp   = "Proxy0x-KillSwitch-AllowApp"
	fwRuleAllowLocal = "Proxy0x-KillSwitch-AllowLocal"
)

func RunElevated(script string) error {
	tmpDir := os.TempDir()
	batFile := filepath.Join(tmpDir, "proxy0x-fw.bat")
	if err := os.WriteFile(batFile, []byte(script), 0o600); err != nil {
		return fmt.Errorf("write script: %w", err)
	}
	defer os.Remove(batFile)

	verb, _ := syscall.UTF16PtrFromString("runas")
	exe, _ := syscall.UTF16PtrFromString("cmd.exe")
	args, _ := syscall.UTF16PtrFromString("/C \"" + batFile + "\"")

	ret, _, _ := shellExecute.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(exe)),
		uintptr(unsafe.Pointer(args)),
		0,
		0, // SW_HIDE
	)
	if ret <= 32 {
		return fmt.Errorf("elevation denied or failed (code %d)", ret)
	}
	return nil
}

func Enable() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get exe path: %w", err)
	}

	script := fmt.Sprintf(`@echo off
netsh advfirewall firewall delete rule name="%s" >nul 2>&1
netsh advfirewall firewall delete rule name="%s" >nul 2>&1
netsh advfirewall firewall delete rule name="%s" >nul 2>&1
netsh advfirewall firewall add rule name="%s" dir=out action=block enable=yes >nul
netsh advfirewall firewall add rule name="%s" dir=out action=allow program="%s" enable=yes >nul
netsh advfirewall firewall add rule name="%s" dir=out action=allow remoteip=127.0.0.1,::1 enable=yes >nul
`,
		fwRuleBlock, fwRuleAllowApp, fwRuleAllowLocal,
		fwRuleBlock, fwRuleAllowApp, exePath, fwRuleAllowLocal,
	)
	return RunElevated(script)
}

func Disable() error {
	script := fmt.Sprintf(`@echo off
netsh advfirewall firewall delete rule name="%s" >nul 2>&1
netsh advfirewall firewall delete rule name="%s" >nul 2>&1
netsh advfirewall firewall delete rule name="%s" >nul 2>&1
`, fwRuleBlock, fwRuleAllowApp, fwRuleAllowLocal)
	return RunElevated(script)
}

func Active() bool {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name="+fwRuleBlock)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return len(output) > 100
}
