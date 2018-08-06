package pkg

import (
	"os/exec"
	"bytes"
	"fmt"
	)

func extractResult(cmd *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "" , fmt.Errorf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if errStr != "" {
		return "", fmt.Errorf("%v", errStr)
	}
	return outStr, nil
}
