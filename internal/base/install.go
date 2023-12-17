//go:build go1.17
// +build go1.17

package base

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GoInstall go get path.
func GoInstall(path ...string) error {
	for _, p := range path {
		if !strings.Contains(p, "@") {
			p += "@latest"
		}
		fmt.Printf("go install %s\n", p)
		cmd := exec.Command("go", "install", p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// Run a command.
func Run(command string, inputs ...string) (reply string, err error) {
	var isVerbose bool
	for _, v := range inputs {
		if v == "verbose" {
			isVerbose = true
		}
	}
	if isVerbose {
		fmt.Printf(fmt.Sprintf("shell: %s %s\n", command, strings.Join(inputs, " ")))
	}
	cmd := exec.Command(command, inputs...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		return
	}
	if stderr.String() != "" {
		err = errors.New(stderr.String())
		return
	}
	reply = stdout.String()
	return
}
