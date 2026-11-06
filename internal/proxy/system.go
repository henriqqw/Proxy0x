package proxy

import (
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

var wininet = syscall.NewLazyDLL("wininet.dll")
var internetSetOption = wininet.NewProc("InternetSetOptionW")

const regPath = `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`

func SetSystemProxy(port int) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	bypass := "localhost;127.*;10.*;192.168.*;<local>"

	runHidden("reg", "add", regPath, "/v", "ProxyServer", "/t", "REG_SZ", "/d", addr, "/f")
	runHidden("reg", "add", regPath, "/v", "ProxyOverride", "/t", "REG_SZ", "/d", bypass, "/f")
	runHidden("reg", "add", regPath, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f")
	refreshProxy()
	return nil
}

func ClearSystemProxy() {
	runHidden("reg", "add", regPath, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f")
	refreshProxy()
}

func runHidden(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	return cmd.Run()
}

func refreshProxy() {
	internetSetOption.Call(0, 39, uintptr(unsafe.Pointer(nil)), 0)
	internetSetOption.Call(0, 37, uintptr(unsafe.Pointer(nil)), 0)
}
