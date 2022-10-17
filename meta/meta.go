package meta

import "fmt"

// GITCOMMIT indicates which git hash the binary was built off of.
var GITCOMMIT string

// VERSION indicates which version of the binary is running.
var VERSION string

// Release holds the current release number, should match the value
// in $GOPATH/src/github.com/kashav/fsql/VERSION.
const Release = "0.5.0"

// Meta returns the version/commit string.
func Meta() string {
	version, commit := VERSION, GITCOMMIT
	if commit == "" || version == "" {
		version, commit = Release, "master"
	}
	return fmt.Sprintf("fsql version %v, built off %v", version, commit)
}
