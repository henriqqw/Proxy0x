package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"proxy_ui/internal/fsutil"
)

type Config struct {
	Language            string `json:"language"`
	Theme               string `json:"theme"`
	AutoRotateInterval  int    `json:"autoRotateInterval"`
	HealthCheckInterval int    `json:"healthCheckInterval"`
}

func DataDir() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "Proxy0x")
}

func DataFile() string {
	return filepath.Join(DataDir(), "proxies.json")
}

func FilePath() string {
	return filepath.Join(DataDir(), "config.json")
}

func Load() Config {
	data, err := os.ReadFile(FilePath())
	if err != nil {
		return defaults()
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaults()
	}
	if cfg.Language == "" {
		cfg.Language = "en"
	}
	if cfg.Theme == "" {
		cfg.Theme = "dark"
	}
	if cfg.HealthCheckInterval <= 0 {
		cfg.HealthCheckInterval = 180
	}
	return cfg
}

func Save(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return fsutil.WriteFileAtomic(FilePath(), data, 0o600)
}

func defaults() Config {
	return Config{
		Language:            "en",
		Theme:               "dark",
		HealthCheckInterval: 180,
	}
}
