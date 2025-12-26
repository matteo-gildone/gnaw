package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const ConfigFileName = ".gnaw"

var (
	ErrEmptyTokensDir    = errors.New("TokensDir is empty")
	ErrEmptySnapshotFile = errors.New("SnapshotFile is empty")
)

type Config struct {
	TokensDir    string `json:"tokensDir"`
	SnapshotFile string `json:"snapshotFile"`
}

func (c Config) Validate() error {
	if c.TokensDir == "" {
		return ErrEmptyTokensDir
	}

	if c.SnapshotFile == "" {
		return ErrEmptySnapshotFile
	}

	return nil
}

func Load(dir string) (Config, error) {
	path := filepath.Join(dir, ConfigFileName)

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if err = cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func Save(dir string, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	path := filepath.Join(dir, ConfigFileName)

	data, err := json.MarshalIndent(cfg, "", "  ")

	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
