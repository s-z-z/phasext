package util

import (
	"io"
	"os/exec"
	"strings"
)

func ExecShell(scriptContent string, stdout, stderr io.Writer) error {
	cmd := exec.Command("bash")
	cmd.Stdin = strings.NewReader(scriptContent)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func ExecCmd(a ...string) error {
	cmd := exec.Command("x", a...)
	return cmd.Run()
}
