package cmd

import (
	"io"
	"fmt"
	"path/filepath"
	"sort"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	usecli "github.com/majestic-fox/ks/pkg/user"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

type namespaceCmd struct {
	configAccess clientcmd.ConfigAccess
	configPath   string
	namespace    string
	out          io.Writer
}

func newNamespaceCmd(out io.Writer) *cobra.Command {
	ns := &namespaceCmd{
		configAccess: clientcmd.NewDefaultPathOptions(),
		out:          out,
	}

	cmd := &cobra.Command{
		Use:     "namespace",
		Aliases: []string{"ns"},
		Short:   "switch current namespace (alias: ns)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if home := homedir.HomeDir(); home != "" {
				ns.configPath = filepath.Join(home, ".kube", "config")
			} else {
				return fmt.Errorf(`there was an error trying to get your kubeconfig ¯\_(ツ)_/¯`)
			}
			return ns.run()
		},
	}
	return cmd
}

func (ns *namespaceCmd) run() (error) {
	config, err := ns.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	kubeCfg, err := clientcmd.BuildConfigFromFlags("", ns.configPath)
	if err != nil {
		return err
	}

	clientSet, err := kubernetes.NewForConfig(kubeCfg)
	if err != nil {
		return err
	}

	namespaceList, err := clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get namespace list: %v", err)
	}

	namespaces := getNamespaceNames(namespaceList)
	sort.Strings(namespaces)

	ns.namespace, err = usecli.Prompt(
		&survey.Select{
			Message:  "Choose a namespace: ",
			Options:  namespaces,
			Default:  config.Contexts[config.CurrentContext].Namespace,
			PageSize: len(namespaces),
		})
	if err != nil {
		return err
	}

	config.Contexts[config.CurrentContext].Namespace = ns.namespace
	err = clientcmd.ModifyConfig(ns.configAccess, *config, true)
	if err != nil {
		return err
	}

	fmt.Fprintf(ns.out, "successfully switched to namespace: %v", ns.namespace)
	return nil
}

func getNamespaceNames(all *corev1.NamespaceList) []string {
	list := make([]string, 0)
	for _, namespace := range all.Items {
		list = append(list, namespace.Name)
	}
	return list
}
