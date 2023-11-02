package cmd

import (
	"errors"
	"fmt"
	"github.com/sergiomss/ks/pkg/user"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"io"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
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
			if len(args) > 0 {
				options.contextName = args[0]
			}
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

	if ctx.contextName == "" {
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

	_, err = getCurrentNamespace(clientcmd.NewDefaultPathOptions())
	if err != nil {
		return fmt.Errorf("failed to get current namespace: %v", err)
	}

	fmt.Fprintf(ctx.out, "Successfully switched to context: %v\n", ctx.contextName)
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
