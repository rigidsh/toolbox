package main

import (
	"fmt"
	"github.com/rigidsh/toolbox/executor"
	"github.com/rigidsh/toolbox/executor/bash"
	"github.com/rigidsh/toolbox/executor/folder"
	"github.com/rigidsh/toolbox/executor/js"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func main() {

	runCmd := filepath.Base(os.Args[0])
	if runCmd == "toolbox" {
		toolbox()
	} else {
		_ = loadExecutor(runCmd).Execute(os.Args[1:])
	}
}

func loadExecutor(name string) executor.Executor {
	baseExecutor := folder.NewFolderExecutor(filepath.Join(installPath, name))
	baseExecutor.RegisterExecutor(".js", func(filePah string) executor.Executor {
		return js.NewExecutor(filePah)
	})
	baseExecutor.RegisterExecutor(".mjs", func(filePah string) executor.Executor {
		return js.NewExecutor(filePah)
	})
	baseExecutor.RegisterExecutor(".sh", func(filePah string) executor.Executor {
		return bash.NewExecutor(filePah)
	})

	return baseExecutor
}

func toolbox() {
	if os.Args[1] == "--complete" {
		baseExecutor := loadExecutor(os.Args[2])
		completeIndex, _ := strconv.Atoi(os.Args[3])
		doAutocomplete(baseExecutor, completeIndex, os.Args[4:])
	} else if os.Args[1] == "--install" {
		err := doInstall(os.Args[2], os.Args[3])
		if err != nil && os.IsPermission(err) {
			args := make([]string, 0, len(os.Args)+1)
			args = append(args, "-E")
			for i, t := range os.Args {
				//TODO: fix it O_O
				if i == 0 {
					tt, _ := which(t)
					args = append(args, tt)
				} else {
					args = append(args, t)
				}
			}
			sudo := exec.Command("sudo", args...)
			sudo.Stdin = os.Stdin
			sudo.Stdout = os.Stdout
			sudo.Stderr = os.Stderr
			_ = sudo.Run()
		}
	}
}

func doAutocomplete(executor executor.Executor, completeIndex int, args []string) {
	var completedArgs []string
	var toComplete string

	if len(args) == completeIndex {
		completedArgs = args[1:]
		toComplete = ""
	} else {
		completedArgs = args[1:completeIndex]
		toComplete = args[completeIndex]
	}

	options, err := executor.Autocomplete(completedArgs, toComplete)
	if err != nil {
	}

	for _, option := range options {
		fmt.Println(option)
	}
}
