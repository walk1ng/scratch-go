package cmd

import "github.com/spf13/cobra"

var multiplyCmd = &cobra.Command{
	Use:   "multiply num1 num2 [numN]",
	Short: "Multiply subcommand add all passed args",

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(multiplyCmd)
}
