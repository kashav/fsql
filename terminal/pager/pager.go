package pager

import (
	"bytes"
	"os"
	"os/exec"
)

const cmd = "less"

var opts = []string{"-S"}

// CommandExists returns true iff cmd exists on the host machine.
func CommandExists() bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// New invokes cmd with in provided as stdin, if cmd is unavailable, return
// an error.
func New(in []byte) error {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}
	pager := exec.Command(path, opts...)
	pager.Stdin = bytes.NewReader(in)
	pager.Stdout = os.Stdout
	pager.Stderr = os.Stderr
	return pager.Run()
}
