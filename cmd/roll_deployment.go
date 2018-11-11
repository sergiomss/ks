package cmd

import (
	"io"
	"fmt"
	"sort"
	"strings"
	"time"
	"os/user"
	"path/filepath"

	usecli "github.com/majestic-fox/ks/pkg/user"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
)

type rollDeployCmd struct {
	configAccess clientcmd.ConfigAccess
	configPath   string
	deployName   string
	out          io.Writer
}

func newRollDeploymentCmd(out io.Writer) *cobra.Command {
	rd := &rollDeployCmd{
		configAccess: clientcmd.NewDefaultPathOptions(),
		out:          out,
	}

	cmd := &cobra.Command{
		Use:     "roll-deployment",
		Aliases: []string{"rd"},
		Short:   "rolls a deployment by adding annotations (alias: rd)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if home := homedir.HomeDir(); home != "" {
				rd.configPath = filepath.Join(home, ".kube", "config")
			} else {
				return fmt.Errorf(`there was an error trying to get your kubeconfig ¯\_(ツ)_/¯`)
			}
			return rd.run()
		},
	}
	return cmd
}

func (rd *rollDeployCmd) run() (error) {
	config, err := rd.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	kubeCfg, err := clientcmd.BuildConfigFromFlags("", rd.configPath)
	if err != nil {
		return err
	}

	clientSet, err := kubernetes.NewForConfig(kubeCfg)
	if err != nil {
		return err
	}

	deploymentsClient := clientSet.AppsV1().Deployments(config.Contexts[config.CurrentContext].Namespace)
	deployList, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment list: %v", err)
	}

	deployments := getDeploymentNames(deployList)
	if len(deployments) == 0 {
		fmt.Fprintln(rd.out, "There are no deployments in the current namespace")
		return nil
	}
	sort.Strings(deployments)

	rd.deployName, err = usecli.Prompt(
		&survey.Select{
			Message:  "Choose a deployment to roll: ",
			Options:  deployments,
			PageSize: len(deployments),
		})
	if err != nil {
		return err
	}

	reason, err := usecli.Prompt(
		&survey.Input{Message: fmt.Sprintf("What is the reason for rolling %v?", rd.deployName)})
	if err != nil {
		return err
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(rd.deployName, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("failed to get latest version of Deployment: %v", getErr))
		}

		user, err := user.Current()
		if err != nil {
			panic(fmt.Errorf("failed to get current user: %v", err))
		}

		annotations := map[string]string{}
		if result.Spec.Template.ObjectMeta.Annotations != nil {
			annotations = result.Spec.Template.ObjectMeta.Annotations
		} else {
			annotations = make(map[string]string)
		}

		annotations["rollout/message"] = strings.TrimSpace(reason)
		annotations["rollout/restarted_at"] = time.Now().UTC().Format(time.RFC3339)
		annotations["rollout/restarted_by"] = user.Name

		result.Spec.Template.ObjectMeta.Annotations = annotations

		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		return fmt.Errorf("update failed: %v", retryErr)
	}
	fmt.Println("Updated deployment...")

	return nil
}

func getDeploymentNames(all *appsv1.DeploymentList) []string {
	list := make([]string, 0)
	for _, deployment := range all.Items {
		list = append(list, deployment.Name)
	}
	return list
}
