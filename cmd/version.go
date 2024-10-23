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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
