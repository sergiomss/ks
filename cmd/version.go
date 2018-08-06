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
		Short: "",
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
	table.AddRow("VERSION", Version)
	table.AddRow("BUILD_DATE", BuildDate)
	fmt.Fprintln(v.out, table)
	return nil
}

