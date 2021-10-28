package cmd

import "github.com/spf13/cobra"

var minusCmd = &cobra.Command{
	Use:   "minus num1 num2 [numN]",
	Short: "Minus subcommand add all passed args",

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(minusCmd)
}
