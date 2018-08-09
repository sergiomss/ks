package cmd

import (
	"io"

	"github.com/spf13/cobra"
	ks "github.com/majestic-fox/ks/pkg"
)

type contextCmd struct {
	out io.Writer
}

func newContextCmd(out io.Writer) *cobra.Command {
	c := &contextCmd{out: out}
	cmd := &cobra.Command{
		Use:   "context",
		Aliases: []string{"c"},
		Short: "switch current context (alias: c)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.run(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func (ns *contextCmd) run() error {
	return ks.SwitchContext()
}
