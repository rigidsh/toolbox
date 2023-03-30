package python

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func createNamedPipeToRead() (name string, err error) {
	name = createPipeRandomName()
	err = syscall.Mkfifo(name, 0666)
	if err != nil {
		return
	}

	return
}

func createPipeRandomName() string {
	return filepath.Join(os.TempDir(), "pythonpipe-"+strconv.Itoa(rand.Int())+".pipe")
}

func callFunction(path string, functionName string, args []interface{}, result *interface{}) (err error) {
	pipeName, err := createNamedPipeToRead()

	if err != nil {
		return
	}

	defer os.Remove(pipeName)

	argsJson, err := json.Marshal(args)
	if err != nil {
		return err
	}

	runScript := fmt.Sprintf(`
import json
from importlib.machinery import SourceFileLoader

module = SourceFileLoader("*", "%s").load_module()
result = module.%s(json.loads('%s'))
pipe = open("%s", "w")
pipe.write(json.dumps(result))
pipe.close()
`,
		path,
		functionName,
		string(argsJson),
		pipeName,
	)

	var buf bytes.Buffer

	cmdExec := exec.Command("python", "-c", runScript)
	cmdExec.Stderr = os.Stderr
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout

	go func() {
		pipe, err := os.OpenFile(pipeName, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			return
		}
		defer pipe.Close()

		tmp := make([]byte, 1024)
		for {
			size, err := pipe.Read(tmp)
			if err == io.EOF {
				return
			}
			buf.Write(tmp[:size])
		}
	}()

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
		return
	}

	return
}
