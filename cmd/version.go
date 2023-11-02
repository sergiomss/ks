package cmd

import (
	"fmt"
	"io"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

type versionCmd struct {
	out io.Writer
}

func newVersionCmd(out io.Writer) *cobra.Command {
	version := &versionCmd{out: out}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get ks's current version and build date",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := version.run(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func (v *versionCmd) run() error {
	table := uitable.New()
	table.AddRow("Version:", Version)
	table.AddRow("BuildDate:", BuildDate)
	table.AddRow("Commit:", Commit)
	fmt.Fprintln(v.out, table)
	return nil
}
