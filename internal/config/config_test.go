package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func ExampleDefault() {
	cfg := Default()
	fmt.Printf("%v", cfg)
	// Output: {tokens tokens.snapshot.json}
}

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.TokenDir == "" {
		t.Error("default TokenDir should be not empty")
	}

	if cfg.SnapshotFile == "" {
		t.Error("default SnapshotFile should be not empty")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr error
	}{
		{
			name: "empty tokens dir",
			cfg: Config{
				TokenDir:     "",
				SnapshotFile: "tokens.snapshot.json"},
			wantErr: ErrEmptyTokenDir,
		},
		{
			name: "empty snapshot file",
			cfg: Config{
				TokenDir:     "tokens",
				SnapshotFile: ""},
			wantErr: ErrEmptySnapshotFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}

		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()

	original := Config{
		TokenDir:     "my-tokens",
		SnapshotFile: "my-snapshot.json",
	}

	if err := Save(tmpDir, original); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	configPath := filepath.Join(tmpDir, ConfigFileName)

	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	loaded, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if loaded.TokenDir != original.TokenDir {
		t.Errorf("TokenDir = %q, want %q", loaded.TokenDir, original.TokenDir)
	}

	if loaded.SnapshotFile != original.SnapshotFile {
		t.Errorf("SnapshotFile = %q, want %q", loaded.SnapshotFile, original.SnapshotFile)
	}
}

func TestLoad_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := Load(tmpDir)

	if err == nil {
		t.Fatal("expected error when config doesn't exists")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestSave_InvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()

	invalid := Config{
		TokenDir:     "",
		SnapshotFile: "my-snapshot.json",
	}

	err := Save(tmpDir, invalid)
	if err == nil {
		t.Fatal("expected error when saving invalid config")
	}
}

func TestLoad_InvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ConfigFileName)
	invalidJSON := `{"tokensDir": "", "snapshotFile": "snapshot.json"}`

	if err := os.WriteFile(configPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	_, err := Load(tmpDir)
	if err == nil {
		t.Fatal("expected error when saving invalid config")
	}
}
