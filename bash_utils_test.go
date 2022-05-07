package main

import (
	"fmt"
	"testing"
)

func TestFindBashPath(t *testing.T) {

	path, _ := which("bash")
	fmt.Println(path)

}
