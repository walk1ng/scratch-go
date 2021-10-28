package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version subcommand show mygit version info.",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := ExecuteCommand("git", "version", args...)
		if err != nil {
			Error(cmd, args, err)
		}

		fmt.Fprint(os.Stdout, out)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
