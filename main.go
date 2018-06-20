package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	survey "gopkg.in/AlecAivazis/survey.v1"
)

func main() {
	curr, err := getCurrentContext()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("current context:", curr)

	ctx, err := getAllContexts()
	if err != nil {
		log.Fatal(err)
	}
	result, err := promptUser(curr, strings.Split(ctx, " "))
	if err != nil {
		log.Fatal(err)
	}

	if result != curr {
		fmt.Println("Selection is different than current context, switching...")
		_, err := changeCurrentContext(result)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Selection matches current context, ignoring...")
	}
}

func promptUser(current string, opts []string) (string, error) {
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

func extractResult(cmd *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if errStr != "" {
		return "", fmt.Errorf("%v", errStr)
	}
	return outStr, nil
}
