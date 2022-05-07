package folder

import (
	"errors"
	"github.com/rigidsh/toolbox/executor"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

type ExecutorFactory func(filePah string) executor.Executor

type Executor struct {
	// file extension to executor factory
	subExecutors map[string]ExecutorFactory
	basePath     string
}

func NewFolderExecutor(basePath string) *Executor {
	return &Executor{
		subExecutors: make(map[string]ExecutorFactory, 0),
		basePath:     basePath,
	}
}

func (folderExecutor *Executor) RegisterExecutor(extension string, executorFactory ExecutorFactory) {
	folderExecutor.subExecutors[extension] = executorFactory
}

func (folderExecutor *Executor) Execute(args []string) error {
	return folderExecutor.execute(folderExecutor.basePath, args)
}

func (folderExecutor *Executor) execute(basePath string, args []string) error {
	if len(args) == 0 {
		return errors.New("unknown argument")
	}

	baseArg := args[0]
	restArgs := args[1:]
	dirItems, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, item := range dirItems {
		if item.IsDir() {
			if item.Name() == baseArg {
				return folderExecutor.execute(filepath.Join(basePath, baseArg), restArgs)
			}
		} else {
			fileExtension := path.Ext(item.Name())
			fileName := strings.TrimSuffix(item.Name(), fileExtension)
			if fileName == baseArg {
				executorFactory, ok := folderExecutor.subExecutors[fileExtension]
				if !ok {
					return errors.New("can't find registered executor")
				}

				return executorFactory(filepath.Join(basePath, item.Name())).Execute(restArgs)
			}
		}
	}

	return errors.New("unknown argument")
}

func (folderExecutor *Executor) Autocomplete(completedArgs []string, toComplete string) ([]string, error) {

	basePath := folderExecutor.basePath

	for i, arg := range completedArgs {
		dirItems, err := ioutil.ReadDir(basePath)
		if err != nil {
			return nil, err
		}

		for _, item := range dirItems {
			if item.IsDir() {
				if item.Name() == arg {
					basePath = filepath.Join(basePath, arg)
					break
				}
			} else {
				fileExtension := path.Ext(item.Name())
				fileName := strings.TrimSuffix(item.Name(), fileExtension)
				if fileName == arg {
					executorFactory, ok := folderExecutor.subExecutors[fileExtension]
					if !ok {
						return nil, errors.New("can't find registered executor")
					}

					return executorFactory(filepath.Join(basePath, item.Name())).Autocomplete(completedArgs[i+1:], toComplete)
				}
			}
		}
	}

	dirItems, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)

	for _, item := range dirItems {
		var argumentSuggestion string

		if item.IsDir() {
			argumentSuggestion = item.Name()
		} else {
			fileExtension := path.Ext(item.Name())
			argumentSuggestion = strings.TrimSuffix(item.Name(), fileExtension)
			_, ok := folderExecutor.subExecutors[fileExtension]
			if !ok {
				continue
			}
		}

		if strings.HasPrefix(argumentSuggestion, toComplete) {
			result = append(result, argumentSuggestion)
		}
	}

	return result, nil

}
