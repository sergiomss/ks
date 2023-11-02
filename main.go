package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"sort"
)

var CLI struct {
	Context struct {
		name struct{} `cmd:"" optional:"" help:"Name of the context to use (optional)"`
	} `cmd:"" help:"Switch kube context" aliases:"ctx,c"`
	Namespace struct {
		name struct{} `cmd:"" optional:"" help:"Name of the namespace to use (optional)"`
	} `cmd:"" help:"Switch kube namespace" aliases:"ns"`
}

func main() {
	ctx := kong.Parse(&CLI)
	if _, err := tea.NewProgram(initialModel(ctx.Command())).Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel(t string) model {
	var options []string

	switch t {
	case "context", "ctx", "c":
		config, err := clientcmd.NewDefaultPathOptions().GetStartingConfig()
		if err != nil {
			panic(err)
		}
		options = contexts(config)
	case "namespace", "ns":
		fmt.Println("namespace")
	default:
		fmt.Sprintf("Command %s not recognized.", t)
		os.Exit(1)
	}

	return model{
		options:    options,
		choiceType: t,
	}
}

func contexts(config *clientcmdapi.Config) []string {
	list := make([]string, 0)
	for name := range config.Contexts {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}
