package cmd

import (
	"log"
	"os"

	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/grpclog"
)

var (
	// Version contains the version of spade
	Version string
	// BuildDate contains the build date of spade
	BuildDate string
)

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ks",
		Short:        "kubeswitch command line tool",
		SilenceUsage: true,
	}

	flags := cmd.PersistentFlags()
	out := cmd.OutOrStdout()

	cmd.AddCommand(
		newVersionCmd(out),
		newContextCmd(out),
		newNamespaceCmd(out),
	)

	flags.Parse(args)
	return cmd
}

func init() {
	// Tell gRPC not to log to console.
	grpclog.SetLogger(log.New(ioutil.Discard, "", log.LstdFlags))
}

func Execute(v, bd string) {
	Version = v
	BuildDate = bd
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func checkArgsLength(argsReceived int, requiredArgs ...string) error {
	expectedNum := len(requiredArgs)
	if argsReceived != expectedNum {
		arg := "arguments"
		if expectedNum == 1 {
			arg = "argument"
		}
		return errors.Errorf("This command needs %v %s: %s", expectedNum, arg, strings.Join(requiredArgs, ", "))
	}
	return nil
}
