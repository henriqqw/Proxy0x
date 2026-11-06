package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

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
