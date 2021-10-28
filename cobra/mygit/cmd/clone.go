package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone repo_url [destination]",
	Short: "clone a repository to destination",

	Run: func(cmd *cobra.Command, args []string) {
		out, err := ExecuteCommand("git", "clone", args...)
		if err != nil {
			Error(cmd, args, err)
		}

		fmt.Fprint(os.Stdout, out)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
