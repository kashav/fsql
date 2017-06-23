package pager

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

var cmds = [...]string{"less", "more"}

// findPagerPath iterates through cmds until finding an executable path on
// the host machine.
func findPagerPath() (path string, err error) {
	for _, cmd := range cmds {
		if path, err = exec.LookPath(cmd); err == nil {
			return path, nil
		}
	}
	return "", errors.New("failed to find command")
}

// New invokes an available pager exectuable with in provided as stdin.
func New(in []byte) error {
	path, err := findPagerPath()
	if err != nil {
		return err
	}
	pager := exec.Command(path)
	pager.Stdin = bytes.NewReader(in)
	pager.Stdout = os.Stdout
	pager.Stderr = os.Stderr
	return pager.Run()
}
