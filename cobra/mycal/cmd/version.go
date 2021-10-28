package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",

	Short: "version subcommand print version info",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Fprintln(os.Stdout, "mycal version 0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
