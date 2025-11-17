package cli

import (
	"github.com/spf13/cobra"
)

type cobraCLI struct {
	root *cobra.Command
}

type cobraCommand struct {
	cmd *cobra.Command
}

func (c *cobraCommand) Run() error {
	return c.cmd.Execute()
}

func NewCobraCLI(root *cobra.Command) CLI {
	return &cobraCLI{root: root}
}

func (c *cobraCLI) AddCommand(name string, cmd Command) {
	if cc, ok := cmd.(*cobraCommand); ok {
		c.root.AddCommand(cc.cmd)
	}
}

func (c *cobraCLI) Run() error {
	return c.root.Execute()
}

func (c *cobraCLI) RunArgs(args []string) error {
	c.root.SetArgs(args)
	return c.root.Execute()
}
