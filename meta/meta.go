package meta

import "fmt"

// GITCOMMIT indicates which git hash the binary was built off of.
var GITCOMMIT string

// VERSION indicates which version of the binary is running.
var VERSION string

const majorRelease = "0.3.x"

// Version returns the version/commit string.
func Version() string {
	version, commit := VERSION, GITCOMMIT
	if commit == "" || version == "" {
		version, commit = majorRelease, "master"
	}
	return fmt.Sprintf("fsql version %v, built off %v", version, commit)
}
