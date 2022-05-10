package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func CallFunction(path string, functionName string, args []interface{}, result *interface{}) error {
	argsJson, err := json.Marshal(args)
	if err != nil {
		return err
	}

	runScript := fmt.Sprintf("import('%s').then(it => it.%s.apply(it, JSON.parse('%s'))).then(it => JSON.stringify(it)).then(it => console.error(it))",
		path,
		functionName,
		string(argsJson),
	)

	var buf bytes.Buffer

	cmdExec := exec.Command("node", "-e", runScript)
	cmdExec.Stderr = &buf
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout

	err = cmdExec.Start()
	if err != nil {
		return err
	}
	err = cmdExec.Wait()
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf.Bytes(), result)
	if err != nil {
		return err
	}

	return nil
}
