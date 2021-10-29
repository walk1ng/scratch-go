package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var minusCmd = &cobra.Command{
	Use:   "minus num1 num2 [numN]",
	Short: "Minus subcommand add all passed args",

	Run: func(cmd *cobra.Command, args []string) {
		nums := convertArgsToFloatSlice(args, ErrorHandling(parseHandling))
		result := calc(nums, MINUS)
		fmt.Fprintf(os.Stdout, "%s = %.2f\n", strings.Join(args, "-"), result)
	},
}

func init() {
	rootCmd.AddCommand(minusCmd)
}
