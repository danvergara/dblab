package cmd

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	// Revision is taken from the vcs.revision tag in Go 1.18+.
	Revision = "unknown"
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, kv := range info.Settings {
		if kv.Value == "" {
			continue
		}

		switch kv.Key {
		case "vcs.revision":
			Revision = kv.Value
		}
	}
}

// NewVersionCmd return a versionCmd instance.
func NewVersionCmd() *cobra.Command {
	// versionCmd represents the version command.
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "The version of the project",
		Long: `The current version of the project.
		This projects follows the semantic versioning standard.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			buildInfo, ok := debug.ReadBuildInfo()
			if !ok {
				return fmt.Errorf("unable to determine version information")
			}

			if buildInfo.Main.Version != "" {
				fmt.Fprintln(cmd.OutOrStdout(), parseVersion(buildInfo.Main.Version))
			} else {
				return errors.New("version: unknown")
			}

			return nil
		},
	}

	return versionCmd
}

// parseVersion parses the version passed as a parameter.
// If the version is equal to unknown or (devel), it shows the commit hash as a revision.
func parseVersion(version string) string {
	if version == "unknown" || version == "(devel)" {
		commit := Revision
		if len(commit) > 7 {
			commit = commit[:7]
		}

		return fmt.Sprintf("rev: %s", commit)
	}

	return version
}

func init() {
	rootCmd.AddCommand(NewVersionCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
