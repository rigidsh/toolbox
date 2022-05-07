package js

import (
	"bufio"
	"bytes"
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

	cmdExec := exec.Command("node", "-e", prepareNodeCommand(executor.filePath, "exec", args))
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

func (executor *Executor) Autocomplete(completedArgs []string, toComplete string) ([]string, error) {

	rawCompletedArgs := make([]string, len(completedArgs)+1)
	rawCompletedArgs[0] = toComplete
	for i, v := range completedArgs {
		rawCompletedArgs[i] = v
	}

	cmdExec := exec.Command("node", "-e", prepareNodeCommand(executor.filePath, "autocomplete", rawCompletedArgs))

	cmdExec.Stderr = os.Stderr
	cmdExec.Stdin = os.Stdin
	var outBuffer bytes.Buffer
	cmdExec.Stdout = bufio.NewWriter(&outBuffer)

	err := cmdExec.Start()
	if err != nil {
		return nil, err
	}
	err = cmdExec.Wait()
	if err != nil {
		return nil, err
	}

	return parseNodeArrayResult(outBuffer.String())
}
