package main

import (
	"fmt"
	"runtime/debug"

	"github.com/danvergara/dblab/cmd"
)

// these values are automagically populated by Goreleaser.
var (
	version  = "dev"
	Revision = "unknown"
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			Revision = kv.Value
		}
	}
}

func main() {
	if version == "dev" {
		version = parseVersion()
	} else {
		// Goreleaser doesn't prefix with a `v`, which we expect.
		version = "v" + version
	}

	cmd.SetVersionInfo(version)
	cmd.Execute()
}

// parseVersion parses the version passed as a parameter.
// If the version is equal to unknown or (devel), it shows the commit hash as a revision.
func parseVersion() string {
	info, _ := debug.ReadBuildInfo()
	v := info.Main.Version

	if v == "unknown" || v == "(devel)" {
		if Revision != "unknown" && Revision != "" {
			commit := Revision
			if len(commit) > 7 {
				commit = commit[:7]
			}
			return fmt.Sprintf("rev: %s", commit)
		}
	} else {
		return v
	}

	return "unknown"
}
