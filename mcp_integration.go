package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type MCPStatus struct {
	Installed   bool   `json:"installed"`
	ExePath     string `json:"exePath"`
	ExeExists   bool   `json:"exeExists"`
	ConfigPath  string `json:"configPath"`
	HealthOK    bool   `json:"healthOk"`
	HealthError string `json:"healthError,omitempty"`
}

func claudeConfigPath() string {
	appdata := os.Getenv("APPDATA")
	if appdata == "" {
		dir, _ := os.UserConfigDir()
		appdata = dir
	}
	return filepath.Join(appdata, "Claude", "claude_desktop_config.json")
}

func mcpExePath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Join(filepath.Dir(exe), "proxy0x-mcp.exe")
}

func (a *App) GetMCPStatus() MCPStatus {
	exePath := mcpExePath()
	configPath := claudeConfigPath()

	status := MCPStatus{
		ExePath:    exePath,
		ConfigPath: configPath,
	}

	if _, err := os.Stat(exePath); err == nil {
		status.ExeExists = true
	}

	data, err := os.ReadFile(configPath)
	if err == nil {
		var cfg map[string]any
		if json.Unmarshal(data, &cfg) == nil {
			if servers, ok := cfg["mcpServers"].(map[string]any); ok {
				if _, ok := servers["proxy0x"]; ok {
					status.Installed = true
				}
			}
		}
	}

	return status
}

func (a *App) InstallMCPConfig() error {
	exePath := mcpExePath()
	if _, err := os.Stat(exePath); err != nil {
		return fmt.Errorf("MCP server binary not found at %s", exePath)
	}

	configPath := claudeConfigPath()
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	var cfg map[string]any
	data, err := os.ReadFile(configPath)
	if err == nil {
		json.Unmarshal(data, &cfg)
	}
	if cfg == nil {
		cfg = make(map[string]any)
	}

	servers, ok := cfg["mcpServers"].(map[string]any)
	if !ok {
		servers = make(map[string]any)
	}

	servers["proxy0x"] = map[string]any{
		"command": exePath,
		"args":    []string{},
	}
	cfg["mcpServers"] = servers

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, out, 0o644)
}

func (a *App) UninstallMCPConfig() error {
	configPath := claudeConfigPath()

	data, err := os.ReadFile(configPath)
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

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, out, 0o644)
}

func (a *App) CheckMCPHealth() MCPStatus {
	status := a.GetMCPStatus()

	if !status.ExeExists {
		status.HealthError = "MCP server binary not found"
		return status
	}

	cmd := exec.Command(status.ExePath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		status.HealthError = fmt.Sprintf("pipe error: %v", err)
		return status
	}

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		status.HealthError = fmt.Sprintf("failed to start: %v", err)
		return status
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
		status.HealthError = "no response from server"
		return status
	}

	for _, line := range bytes.Split(stdout.Bytes(), []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var resp map[string]any
		if json.Unmarshal(line, &resp) == nil {
			if _, ok := resp["result"]; ok {
				status.HealthOK = true
				return status
			}
		}
	}

	status.HealthError = "invalid response from server"
	return status
}
