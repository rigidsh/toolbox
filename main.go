package main

import (
	"fmt"
	"github.com/rigidsh/toolbox/executor"
	"github.com/rigidsh/toolbox/executor/bash"
	"github.com/rigidsh/toolbox/executor/folder"
	"github.com/rigidsh/toolbox/executor/js"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const installPath = "/var/toolbox"

func main() {
	f, err := os.OpenFile("complete.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	log.SetOutput(f)

	runCmd := filepath.Base(os.Args[0])
	if runCmd == "toolbox" {
		toolbox()
	} else {
		err := loadExecutor(runCmd).Execute(os.Args[1:])
		if err != nil {
			log.Println(err)
		}
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
	log.Printf("all args: %v", os.Args)
	if os.Args[1] == "--complete" {
		baseExecutor := loadExecutor(os.Args[2])
		completeIndex, _ := strconv.Atoi(os.Args[3])
		doAutocomplete(baseExecutor, completeIndex, os.Args[4:])
	} else if os.Args[1] == "--install" {
		err := doInstall(os.Args[2], os.Args[3])
		if err != nil && os.IsPermission(err) {
			log.Println("Try to use sudo")
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

func doInstall(name, path string) error {
	cmd := filepath.Join("/usr/local/bin/", name)
	var toolboxAbsPath string
	if t, ok := which(os.Args[0]); ok {
		toolboxAbsPath = t
	} else {
		toolboxAbsPath = os.Args[0]
	}
	err := os.Symlink(toolboxAbsPath, cmd)
	if err != nil {
		if os.IsExist(err) {
			err = os.Remove(cmd)
			if err != nil {
				return err
			}
			return doInstall(name, path)
		} else {
			return err
		}
	}
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		err := os.Mkdir(installPath, 0755)
		if err != nil {
			return err
		}
	}
	err = cp(path, filepath.Join(installPath, name))

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		filepath.Join("/usr/share/bash-completion/completions", name),
		[]byte(
			"__"+name+"() { \n"+
				"args=$( IFS=$' '; echo \"${COMP_WORDS[*]}\"); \n"+
				"args=\"--complete "+name+"  ${COMP_CWORD} $args\"; \n"+
				"for i in `"+toolboxAbsPath+" $args`; do COMPREPLY+=($i); done; \n"+
				"}\n"+
				"\n"+
				"complete -F __"+name+" "+name),
		644)
	if err != nil {
		return err
	}

	return nil
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
	log.Printf("Complete args: %v", args)
	log.Printf("Complete index: %d", completeIndex)
	log.Printf("Complete result: %v", options)

	for _, option := range options {
		fmt.Println(option)
	}
}
