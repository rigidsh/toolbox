package bash

import (
	"os"
	"os/exec"
)

func NewExecutor(filePath string) *Executor {
	return &Executor{
		filePath: filePath,
	}
}

type Executor struct {
	filePath string
}

func (executor *Executor) Execute(args []string) error {

	cmdExec := exec.Command("/bin/bash", append([]string{executor.filePath}, args...)...)
	cmdExec.Stderr = os.Stderr
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout

	err := cmdExec.Start()
	if err != nil {
		return err
	}
	err = cmdExec.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (executor *Executor) Autocomplete(_ []string, _ string) ([]string, error) {
	return []string{}, nil
}
