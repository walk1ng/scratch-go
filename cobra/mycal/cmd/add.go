package cmd

import "github.com/spf13/cobra"

var addCmd = &cobra.Command{
	Use:   "add num1 num2 [numN]",
	Short: "Add subcommand add all passed args",

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
