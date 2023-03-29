package main

import (
	"os"
	"path/filepath"
)

const installPath = "/var/toolbox"
const binPath = "/usr/local/bin/"
const bashCompletionScriptPath = "/usr/share/bash-completion/completions"

func doInstall(name, path string) error {
	cmd := filepath.Join(binPath, name)
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

	err = createBashCompletionPath(name, toolboxAbsPath)

	if err != nil {
		return err
	}

	return nil
}

func createBashCompletionPath(name, toolboxAbsPath string) error {
	return os.WriteFile(
		filepath.Join(bashCompletionScriptPath, name),
		[]byte(
			"__"+name+"() { \n"+
				"args=$( IFS=$' '; echo \"${COMP_WORDS[*]}\"); \n"+
				"args=\"--complete "+name+"  ${COMP_CWORD} $args\"; \n"+
				"for i in `"+toolboxAbsPath+" $args`; do COMPREPLY+=($i); done; \n"+
				"}\n"+
				"\n"+
				"complete -F __"+name+" "+name),
		644)
}
