package version

import (
	"fmt"
)

// The following variables are overwritten at build time
var (
	Version     = ""
	GoVersion   = ""
	Commit      = ""
	CompileTime = ""
)

// FullVersion formats the version to be printed
func FullVersion() string {
	return fmt.Sprintf("%s, build %s. Built %s using %s", Version, Commit, CompileTime, GoVersion)
}
