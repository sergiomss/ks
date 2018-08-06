package pkg

import (
	"fmt"
	"os/exec"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)
func SwitchContext() error {
	curr, err := getCurrentContext()
	if err != nil {
		return fmt.Errorf("failed to switch context: %v", err)
	}
	fmt.Println("current context:", curr)

	ctx, err := getAllContexts()
	if err != nil {
		return fmt.Errorf("failed to switch context: %v", err)
	}
	result, err := promptUserContext(curr, strings.Split(ctx, " "))
	if err != nil {
		return fmt.Errorf("failed to switch context: %v", err)
	}

	if result != curr {
		fmt.Println("Selection is different than current context, switching...")
		_, err := changeCurrentContext(result)
		if err != nil {
			return fmt.Errorf("failed to switch context: %v", err)
		}
	} else {
		fmt.Println("Selection matches current context, ignoring...")
	}
	return nil
}

func promptUserContext(current string, opts []string) (string, error) {
	var qs = []*survey.Question{
		{
			Name: "context",
			Prompt: &survey.Select{
				Message: "Choose a kubecontext:",
				Options: opts,
				Default: current,
			},
		},
	}

	answers := struct {
		Context string `survey:"context"`
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return "", err
	}
	return answers.Context, nil
}

func changeCurrentContext(ctx string) (string, error) {
	return extractResult(exec.Command("kubectl", "config", "use-context", ctx))
}

func getCurrentContext() (string, error) {
	return extractResult(exec.Command("kubectl", "config", "view", "-o", "jsonpath={.current-context}"))
}

func getAllContexts() (string, error) {
	return extractResult(exec.Command("kubectl", "config", "view", "-o", "jsonpath={.contexts[*].name}"))
}

