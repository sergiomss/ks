package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"io"
	"path/filepath"
	"sort"

	usecli "github.com/sergiomss/ks/pkg/user"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

	currentNs, err := getCurrentNamespace(ns.configAccess)
	if err != nil {
		return fmt.Errorf("failed to get current namespace: %v", err)
	}
	namespaces := getNamespaceNames(namespaceList)
	sort.Strings(namespaces)

	ns.namespace, err = usecli.Prompt(
		&survey.Select{
			Message:  "Choose a namespace: ",
			Options:  namespaces,
			Default:  currentNs,
			PageSize: len(namespaces),
		})
	if err != nil {
		return err
	}

	err = setCurrentNamespace(ns.configAccess, ns.namespace)
	if err != nil {
		return fmt.Errorf("failed to set current namespace to %v: %v", ns.namespace, err)
	}

	fmt.Fprintf(ns.out, "Successfully switched to namespace: %v\n", ns.namespace)
	return nil
}

func getCurrentNamespace(access clientcmd.ConfigAccess) (string, error) {
	config, err := access.GetStartingConfig()
	if err != nil {
		return "", err
	}
	return config.Contexts[config.CurrentContext].Namespace, nil
}

func setCurrentNamespace(access clientcmd.ConfigAccess, namespace string) error {
	config, err := access.GetStartingConfig()
	if err != nil {
		return err
	}
	config.Contexts[config.CurrentContext].Namespace = namespace
	err = clientcmd.ModifyConfig(access, *config, true)
	if err != nil {
		return err
	}
	if err := setTillerNamespace(namespace); err != nil {
		return fmt.Errorf("failed to set tiller namespace: %v", err)
	}
	return nil
}

func setTillerNamespace(namespace string) error {
	return clipboard.WriteAll(fmt.Sprintf("export TILLER_NAMESPACE=%v", namespace))
}

func getNamespaceNames(all *corev1.NamespaceList) []string {
	list := make([]string, 0)
	for _, namespace := range all.Items {
		list = append(list, namespace.Name)
	}
	return list
}
