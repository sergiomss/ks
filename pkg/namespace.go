package pkg

import (
	"fmt"
	"strings"
	"gopkg.in/AlecAivazis/survey.v1"
	"os/exec"
)

func SwitchNamespace() error {
	ctx, err := getCurrentContext()
	if err != nil {
		return fmt.Errorf("failed to get current namespace: %v", err)
	}

	curr, err := getCurrentNamespace(ctx)
	if err != nil {
		return fmt.Errorf("failed to switch namespace: %v", err)
	}
	fmt.Println("current namespace:", curr)

	namespaces, err := getAllNamespaces()
	if err != nil {
		return fmt.Errorf("failed to switch namespace: %v", err)
	}
	result, err := promptUserNamespace(curr, strings.Split(namespaces, " "))
	if err != nil {
		return fmt.Errorf("failed to switch namespace: %v", err)
	}

	if result != curr {
		fmt.Println("Selection is different than current namespace, switching...")
		_, err := changeCurrentNamespace(ctx, result)
		if err != nil {
			return fmt.Errorf("failed to switch namespace: %v", err)
		}
	} else {
		fmt.Println("Selection matches current namespace, ignoring...")
	}
	return nil
}

func promptUserNamespace(current string, opts []string) (string, error) {
	var qs = []*survey.Question{
		{
			Name: "namespace",
			Prompt: &survey.Select{
				Message: "Choose a namespace:",
				Options: opts,
				Default: current,
			},
		},
	}

	answers := struct {
		Context string `survey:"namespace"`
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return "", err
	}
	return answers.Context, nil
}

func changeCurrentNamespace(ctx, ns string) (string, error) {
	namespace := fmt.Sprintf("--namespace=%v", ns)
	return extractResult(exec.Command("kubectl", "config", "set-context", ctx, namespace))
}

func getCurrentNamespace(context string) (string, error) {
	jp := fmt.Sprintf("jsonpath='{.contexts[?(@.name == %q)].context.namespace}'", context)
	return extractResult(exec.Command("kubectl", "config", "view", "-o", jp))
}

func getAllNamespaces() (string, error) {
	return extractResult(exec.Command("kubectl", "get", "namespace", "-o", "jsonpath={.items..name}"))
}
