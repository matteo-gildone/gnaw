package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/matteo-gildone/gnaw/internal/config"
)

var ErrConfigExists = errors.New("configuration already exists")

var cmdInit = &Command{
	UsageLine: "init [-name <file name>] [-tokens <dir>]",
	Short:     "initialise gnaw configuration",
	Long: `Init creates a .gnaw configuration file in the current directory.

Examples:
	gnaw init
	gnaw init -tokens design-tokens -name tokens.snap.json
`,
}

var (
	initTokensDir    string
	initSnapshotName string
)

func init() {
	cmdInit.Run = runInit
	cmdInit.Flag.StringVar(&initSnapshotName, "name", "tokens.snapshot.json", "snapshot file name")
	cmdInit.Flag.StringVar(&initTokensDir, "tokens", "tokens", "tokens directory")
}

func runInit(ctx context.Context, args []string) error {
	if _, err := config.Load("."); err == nil {
		return ErrConfigExists
	}

	cfg := config.Config{
		TokensDir:    initTokensDir,
		SnapshotFile: initSnapshotName,
	}

	if err := config.Save(".", cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Created .gnaw")

	return nil
}
