package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		os.Setenv(k, v.Value)
	}

	currentEnv := os.Environ()
	res := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	res.Env = append([]string{}, currentEnv...)

	if res.Stdin == nil {
		res.Stdin = os.Stdin
	}

	if res.Stdout == nil {
		res.Stdout = os.Stdout
	}

	if res.Stderr == nil {
		res.Stderr = os.Stderr
	}

	err := res.Run()
	if err != nil {
		log.Fatal(err)
	}

	return res.ProcessState.ExitCode()
}
