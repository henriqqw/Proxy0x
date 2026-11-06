package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if os.Getenv("PROXY0X_FAKE_MCP") == "1" {
		fakeMCPServer()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// fakeMCPServer behaves like servers that stop on stdin EOF before replying:
// it answers initialize only if stdin is still open shortly after the request.
func fakeMCPServer() {
	r := bufio.NewReader(os.Stdin)
	if _, err := r.ReadString('\n'); err != nil {
		return
	}
	eof := make(chan struct{})
	go func() {
		if _, err := r.ReadByte(); err != nil {
			close(eof)
		}
	}()
	select {
	case <-eof:
		return
	case <-time.After(200 * time.Millisecond):
		os.Stdout.WriteString(`{"jsonrpc":"2.0","id":1,"result":{}}` + "\n")
	}
	time.Sleep(time.Second)
}

func TestCheckMCPHealthWaitsForTheResponse(t *testing.T) {
	t.Setenv("PROXY0X_FAKE_MCP", "1")
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if err := checkMCPHealth(exe); err != nil {
		t.Fatalf("healthy server reported unhealthy: %v", err)
	}
}

func TestAddProxy0xRefusesToOverwriteInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "claude_desktop_config.json")
	original := `{"preferences": {"theme": "dark"}, ` // truncated mid-write
	os.WriteFile(path, []byte(original), 0o644)

	if err := addProxy0x(path, `C:\x\proxy0x-mcp.exe`); err == nil {
		t.Fatal("expected an error for an unparseable config")
	}
	if got, _ := os.ReadFile(path); string(got) != original {
		t.Fatalf("config was overwritten: %s", got)
	}
}

func TestAddAndRemoveProxy0xPreserveOtherKeysWithBackup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "claude_desktop_config.json")
	original := `{"preferences":{"theme":"dark"},"mcpServers":{"other":{"command":"o.exe"}}}`
	os.WriteFile(path, []byte(original), 0o644)

	if err := addProxy0x(path, `C:\x\proxy0x-mcp.exe`); err != nil {
		t.Fatal(err)
	}
	cfg, err := readJSON(path)
	if err != nil {
		t.Fatal(err)
	}
	servers := cfg["mcpServers"].(map[string]any)
	if cfg["preferences"] == nil || servers["other"] == nil || servers["proxy0x"] == nil {
		t.Fatalf("keys lost or proxy0x missing: %v", cfg)
	}
	if bak, _ := os.ReadFile(path + ".proxy0x.bak"); string(bak) != original {
		t.Fatalf("backup = %q, want the original file", bak)
	}

	if err := removeProxy0x(path); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	var after map[string]any
	json.Unmarshal(data, &after)
	if after["preferences"] == nil || hasProxy0x(after) {
		t.Fatalf("remove lost keys or kept proxy0x: %s", data)
	}
}
