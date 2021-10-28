package cmd

import "github.com/spf13/cobra"

var divideCmd = &cobra.Command{
	Use:   "divide num1 num2 [numN]",
	Short: "Divide subcommand add all passed args",

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(divideCmd)
}
