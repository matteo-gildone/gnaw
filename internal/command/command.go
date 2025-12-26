package command

import (
	"context"
	"flag"
)

type Command struct {
	Run   func(ctx context.Context, args []string) error
	Usage string
	Short string
	Long  string
	Flag  flag.FlagSet
}

func (c *Command) Name() string {
	name := c.Usage
	for i, r := range name {
		if r == ' ' || r == '[' {
			return name[:i]
		}
	}
	return name
}
