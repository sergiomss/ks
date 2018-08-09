package cmd

import (
	"io"

	"github.com/spf13/cobra"
	ks "github.com/majestic-fox/ks/pkg"
	"fmt"
)

type currentNamespaceCmd struct {
	out io.Writer
}

func newCurrentNamespaceCmd(out io.Writer) *cobra.Command {
	cns := &currentNamespaceCmd{out: out}
	cmd := &cobra.Command{
		Use:   "current-namespace",
		Aliases: []string{"current-ns", "cns"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cns.run(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func (cns *currentNamespaceCmd) run() error {
	ctx, err := ks.GetCurrentContext()
	if err != nil {
		fmt.Print("☠️")
	}
	ns, err := ks.GetCurrentNamespace(ctx)
	if err != nil {
		fmt.Print("☠️")
	}
	fmt.Print(ns)
	return nil
}

