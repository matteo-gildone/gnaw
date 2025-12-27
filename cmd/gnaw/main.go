package main

import (
	"context"
	"flag"
	"fmt"
	"os"
)

type Command struct {
	Run       func(ctx context.Context, args []string) error
	UsageLine string
	Short     string
	Long      string
	Flag      flag.FlagSet
}

func (c *Command) Name() string {
	name := c.UsageLine
	for i, r := range name {
		if r == ' ' || r == '[' {
			return name[:i]
		}
	}
	return name
}

var commands = []*Command{
	cmdInit,
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		usage()
		os.Exit(2)
	}

	cmdName := args[0]

	for _, cmd := range commands {
		if cmd.Name() == cmdName {
			cmd.Flag.Usage = func() { usageCommand(cmd) }
			cmd.Flag.Parse(args[1:])
			ctx := context.Background()
			if err := cmd.Run(ctx, cmd.Flag.Args()); err != nil {
				fmt.Fprintf(os.Stderr, "gnaw %s: %v\n", cmdName, err)
				os.Exit(1)
			}
			return
		}
	}
	fmt.Fprintf(os.Stderr, "gnaw unknown command %q\n", cmdName)
	fmt.Fprintf(os.Stderr, "run 'gnaw help' for usage\n")
	os.Exit(2)
}

func usage() {
	fmt.Fprintf(os.Stderr, `gnaw is a tool for managing design tokens snapshot.

Usage:
	gnaw <command> [arguments]

The commands are:

`)

	for _, cmd := range commands {
		fmt.Fprintf(os.Stderr, "    %-12s %s\n", cmd.Name(), cmd.Short)
	}
	fmt.Fprintf(os.Stderr, "\nUse 'gnaw <command> -h for more information about a command.'\n")
}

func usageCommand(cmd *Command) {
	fmt.Fprintf(os.Stderr, "usage: gnaw %s\n\n", cmd.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", cmd.Long)
	cmd.Flag.PrintDefaults()
}
