package js

import (
	"github.com/rigidsh/toolbox/executor/js/node"
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

	var result interface{}
	rawArgs := make([]interface{}, 0, len(args))
	for _, arg := range args {
		rawArgs = append(rawArgs, arg)
	}

	err := node.CallFunction(executor.filePath, "exec", rawArgs, &result)
	if err != nil {
		return err
	}

	return nil
}

func (executor *Executor) Autocomplete(completedArgs []string, toComplete string) ([]string, error) {

	var rawResult interface{}

	rawCompletedArgs := make([]interface{}, len(completedArgs)+1)
	rawCompletedArgs[len(completedArgs)] = toComplete
	for i, v := range completedArgs {
		rawCompletedArgs[i] = v
	}

	err := node.CallFunction(executor.filePath, "autocomplete", rawCompletedArgs, &rawResult)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(rawResult.([]interface{})))
	for _, r := range rawResult.([]interface{}) {
		result = append(result, r.(string))
	}

	return result, nil
}
