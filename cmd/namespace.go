package cmd

import (
	"io"

	"github.com/spf13/cobra"
	ks "github.com/majestic-fox/ks/pkg"
)

type namespaceCmd struct {
	out io.Writer
}

func newNamespaceCmd(out io.Writer) *cobra.Command {
	ns := &namespaceCmd{out: out}
	cmd := &cobra.Command{
		Use:   "namespace",
		Aliases: []string{"ns"},
		Short: "switch current namespace (alias: ns)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ns.run(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func (ns *namespaceCmd) run() error {
	return ks.SwitchNamespace()
}
