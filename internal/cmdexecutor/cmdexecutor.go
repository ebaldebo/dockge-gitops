package cmdexecutor

import (
	"fmt"
	"os/exec"
)

//go:generate moq -rm -out cmdexecutor_mock.go . CommandExecutor
type CommandExecutor interface {
	ExecuteCommand(name string, args ...string) ([]byte, error)
}

type DefaultCommandExecutor struct{}

func (c DefaultCommandExecutor) ExecuteCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("command execution failed: %v, output: %s", err, output)
	}
	return output, nil
}
