package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SetVersionInfo(version string) {
	rootCmd.Version = fmt.Sprint(version)
}

// NewVersionCmd return a versionCmd instance.
func NewVersionCmd() *cobra.Command {
	// versionCmd represents the version command.
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "The version of the project",
		Long: `The current version of the project.
		This projects follows the semantic versioning standard.`,
		Run: func(cmd *cobra.Command, args []string) {
			root := cmd.Root()
			root.SetArgs([]string{"--version"})
			_ = root.Execute()
		},
	}

	return versionCmd
}

func init() {
	rootCmd.AddCommand(NewVersionCmd())
}
