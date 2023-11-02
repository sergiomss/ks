package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
)

var (
	// Version contains the version of ks
	Version string
	// BuildDate contains the build date of ks
	BuildDate string
	// Commit contains the commit hash of ks
	Commit string
)

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ks",
		Short:        "all around dev-ops tool",
		SilenceUsage: true,
	}

	flags := cmd.PersistentFlags()
	out := cmd.OutOrStdout()

	cmd.AddCommand(
		newVersionCmd(out),
		newContextCmd(out),
		newNamespaceCmd(out),
		newRollDeploymentCmd(out),
	)

	flags.Parse(args)
	return cmd
}

func init() {
	// Tell gRPC not to log to console.
	grpclog.SetLogger(log.New(ioutil.Discard, "", log.LstdFlags))
}

func Execute(version, buildDate, commit string) {
	Version = version
	BuildDate = buildDate
	Commit = commit
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
