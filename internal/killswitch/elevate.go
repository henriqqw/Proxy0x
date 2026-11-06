package killswitch

import (
	"errors"
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

// ErrCancelled means the user declined the UAC prompt; nothing was run.
var ErrCancelled = errors.New("administrator permission was not granted")

const (
	seeMaskNoCloseProcess = 0x00000040
	seeMaskNoAsync        = 0x00000100
	seeMaskFlagNoUI       = 0x00000400
	swHide                = 0

	elevatedTimeoutMs = 60_000
)

var procShellExecuteExW = windows.NewLazySystemDLL("shell32.dll").NewProc("ShellExecuteExW")

// shellExecuteInfo mirrors SHELLEXECUTEINFOW.
type shellExecuteInfo struct {
	cbSize         uint32
	fMask          uint32
	hwnd           uintptr
	lpVerb         *uint16
	lpFile         *uint16
	lpParameters   *uint16
	lpDirectory    *uint16
	nShow          int32
	hInstApp       uintptr
	lpIDList       uintptr
	lpClass        *uint16
	hkeyClass      uintptr
	dwHotKey       uint32
	hIconOrMonitor uintptr
	hProcess       windows.Handle
}

// RunElevated runs one cmd.exe command line as administrator (UAC prompt),
// waits for it and fails on a non-zero exit code. The command is passed on
// the command line, never through a file another process could swap.
func RunElevated(command string) error {
	sysDir, err := windows.GetSystemDirectory()
	if err != nil {
		return fmt.Errorf("locate system directory: %w", err)
	}
	verb, err := windows.UTF16PtrFromString("runas")
	if err != nil {
		return err
	}
	file, err := windows.UTF16PtrFromString(filepath.Join(sysDir, "cmd.exe"))
	if err != nil {
		return err
	}
	// /S: cmd strips only the outer quotes and keeps the inner ones as written.
	params, err := windows.UTF16PtrFromString(`/S /C "` + command + `"`)
	if err != nil {
		return err
	}

	info := shellExecuteInfo{
		fMask:        seeMaskNoCloseProcess | seeMaskNoAsync | seeMaskFlagNoUI,
		lpVerb:       verb,
		lpFile:       file,
		lpParameters: params,
		nShow:        swHide,
	}
	info.cbSize = uint32(unsafe.Sizeof(info))

	ok, _, callErr := procShellExecuteExW.Call(uintptr(unsafe.Pointer(&info)))
	if ok == 0 {
		if errors.Is(callErr, windows.ERROR_CANCELLED) {
			return ErrCancelled
		}
		return fmt.Errorf("start elevated command: %w", callErr)
	}
	if info.hProcess == 0 {
		return errors.New("start elevated command: no process handle")
	}
	defer windows.CloseHandle(info.hProcess)

	event, err := windows.WaitForSingleObject(info.hProcess, elevatedTimeoutMs)
	if err != nil {
		return fmt.Errorf("wait for elevated command: %w", err)
	}
	if event != windows.WAIT_OBJECT_0 {
		return errors.New("elevated command timed out")
	}
	var code uint32
	if err := windows.GetExitCodeProcess(info.hProcess, &code); err != nil {
		return fmt.Errorf("read elevated command exit code: %w", err)
	}
	if code != 0 {
		return fmt.Errorf("elevated command failed (exit code %d)", code)
	}
	return nil
}
