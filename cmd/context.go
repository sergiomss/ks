package cmd

import (
	"errors"
	"fmt"
	"io"
	"github.com/majestic-fox/ks/pkg/user"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sort"
)

type contextCmd struct {
	configAccess clientcmd.ConfigAccess
	contextName  string
	out          io.Writer
}

func newContextCmd(out io.Writer) *cobra.Command {
	options := &contextCmd{
		configAccess: clientcmd.NewDefaultPathOptions(),
		out:          out,
	}

	cmd := &cobra.Command{
		Use:     "context",
		Aliases: []string{"c"},
		Short:   "switch current context (alias: c)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.run()
		},
	}
	return cmd
}

func (ctx *contextCmd) run() error {

	config, err := ctx.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	ctxs := contexts(config)
	ctx.contextName, err = user.Prompt(
		&survey.Select{
			Message:  "Choose a context: ",
			Options:  ctxs,
			Default:  config.CurrentContext,
			PageSize: len(ctxs),
		})
	if err != nil {
		return err
	}

	err = ctx.validate(config)
	if err != nil {
		return err
	}

	config.CurrentContext = ctx.contextName
	err = clientcmd.ModifyConfig(ctx.configAccess, *config, true)
	if err != nil {
		return err
	}

	fmt.Fprintf(ctx.out, "successfully switched to context: %v", ctx.contextName)
	return nil
}

func (ctx *contextCmd) validate(config *clientcmdapi.Config) error {
	if len(ctx.contextName) == 0 {
		return errors.New("empty context names are not allowed")
	}

	for name := range config.Contexts {
		if name == ctx.contextName {
			return nil
		}
	}

	return fmt.Errorf("no context exists with the name: %q", ctx.contextName)
}

func contexts(config *clientcmdapi.Config) []string {
	list := make([]string, 0)
	for name := range config.Contexts {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}
