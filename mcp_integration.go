package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"proxy_ui/internal/fsutil"
)

type MCPClientInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Detected   bool   `json:"detected"`
	Installed  bool   `json:"installed"`
	ConfigPath string `json:"configPath"`
}

type MCPStatus struct {
	ExePath     string          `json:"exePath"`
	ExeExists   bool            `json:"exeExists"`
	Clients     []MCPClientInfo `json:"clients"`
	HealthOK    bool            `json:"healthOk"`
	HealthError string          `json:"healthError,omitempty"`
}

type mcpClientDef struct {
	id         string
	name       string
	configPath func() string
	detectDirs func() []string
}

func roamingDir() string {
	if d := os.Getenv("APPDATA"); d != "" {
		return d
	}
	d, _ := os.UserConfigDir()
	return d
}

func homeDir() string {
	d, _ := os.UserHomeDir()
	return d
}

var mcpClientDefs = []mcpClientDef{
	{
		id:   "claude-desktop",
		name: "Claude Desktop",
		configPath: func() string {
			return filepath.Join(roamingDir(), "Claude", "claude_desktop_config.json")
		},
		detectDirs: func() []string {
			return []string{filepath.Join(roamingDir(), "Claude")}
		},
	},
	{
		id:   "codex",
		name: "Codex",
		configPath: func() string {
			return filepath.Join(homeDir(), ".codex", "config.json")
		},
		detectDirs: func() []string {
			return []string{filepath.Join(homeDir(), ".codex")}
		},
	},
	{
		id:   "cursor",
		name: "Cursor",
		configPath: func() string {
			return filepath.Join(homeDir(), ".cursor", "mcp.json")
		},
		detectDirs: func() []string {
			return []string{
				filepath.Join(roamingDir(), "Cursor"),
				filepath.Join(homeDir(), ".cursor"),
			}
		},
	},
	{
		id:   "windsurf",
		name: "Windsurf",
		configPath: func() string {
			return filepath.Join(homeDir(), ".codeium", "windsurf", "mcp_config.json")
		},
		detectDirs: func() []string {
			return []string{
				filepath.Join(roamingDir(), "Windsurf"),
				filepath.Join(homeDir(), ".codeium", "windsurf"),
			}
		},
	},
}

func mcpExePath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Join(filepath.Dir(exe), "proxy0x-mcp.exe")
}

func clientDetected(def mcpClientDef) bool {
	for _, dir := range def.detectDirs() {
		if _, err := os.Stat(dir); err == nil {
			return true
		}
	}
	return false
}

func readJSON(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func hasProxy0x(cfg map[string]any) bool {
	servers, ok := cfg["mcpServers"].(map[string]any)
	if !ok {
		return false
	}
	_, ok = servers["proxy0x"]
	return ok
}

// writeClientConfig never leaves another app's config truncated and keeps
// the previous version next to it.
func writeClientConfig(path string, cfg map[string]any) error {
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return fsutil.WriteFileAtomicWithBackup(path, path+".proxy0x.bak", out, 0o644)
}

func addProxy0x(path, exePath string) error {
	var cfg map[string]any
	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		if err := json.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("%s is not valid JSON, refusing to overwrite it: %w", path, err)
		}
	case !os.IsNotExist(err):
		return err
	}
	if cfg == nil {
		cfg = make(map[string]any)
	}

	servers, _ := cfg["mcpServers"].(map[string]any)
	if servers == nil {
		servers = make(map[string]any)
	}
	servers["proxy0x"] = map[string]any{
		"command": exePath,
		"args":    []string{},
	}
	cfg["mcpServers"] = servers
	return writeClientConfig(path, cfg)
}

func removeProxy0x(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	servers, ok := cfg["mcpServers"].(map[string]any)
	if !ok {
		return nil
	}
	delete(servers, "proxy0x")
	cfg["mcpServers"] = servers
	return writeClientConfig(path, cfg)
}

func (a *App) GetMCPStatus() MCPStatus {
	exePath := mcpExePath()
	status := MCPStatus{ExePath: exePath}
	if _, err := os.Stat(exePath); err == nil {
		status.ExeExists = true
	}
	for _, def := range mcpClientDefs {
		ci := MCPClientInfo{
			ID:         def.id,
			Name:       def.name,
			ConfigPath: def.configPath(),
			Detected:   clientDetected(def),
		}
		if cfg, err := readJSON(ci.ConfigPath); err == nil {
			ci.Installed = hasProxy0x(cfg)
		}
		status.Clients = append(status.Clients, ci)
	}
	return status
}

func (a *App) InstallMCPClient(clientID string) error {
	exePath := mcpExePath()
	if _, err := os.Stat(exePath); err != nil {
		return fmt.Errorf("MCP server binary not found at %s", exePath)
	}
	for _, def := range mcpClientDefs {
		if def.id == clientID {
			return addProxy0x(def.configPath(), exePath)
		}
	}
	return fmt.Errorf("unknown client: %s", clientID)
}

func (a *App) UninstallMCPClient(clientID string) error {
	for _, def := range mcpClientDefs {
		if def.id == clientID {
			return removeProxy0x(def.configPath())
		}
	}
	return fmt.Errorf("unknown client: %s", clientID)
}

func (a *App) CheckMCPHealth() MCPStatus {
	status := a.GetMCPStatus()
	if !status.ExeExists {
		status.HealthError = "MCP server binary not found"
		return status
	}
	if err := checkMCPHealth(status.ExePath); err != nil {
		status.HealthError = err.Error()
		return status
	}
	status.HealthOK = true
	return status
}

// checkMCPHealth starts the MCP server and expects a JSON-RPC initialize result.
func checkMCPHealth(exePath string) error {
	cmd := exec.Command(exePath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("pipe error: %v", err)
	}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start: %v", err)
	}

	initReq := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"proxy0x-healthcheck","version":"1.0.0"}}}` + "\n"
	stdin.Write([]byte(initReq))
	stdin.Close()

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		cmd.Process.Kill()
		<-done
	case <-done:
	}

	if stdout.Len() == 0 {
		return errors.New("no response from server")
	}

	for _, line := range bytes.Split(stdout.Bytes(), []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var resp map[string]any
		if json.Unmarshal(line, &resp) == nil {
			if _, ok := resp["result"]; ok {
				return nil
			}
		}
	}
	return errors.New("invalid response from server")
}
