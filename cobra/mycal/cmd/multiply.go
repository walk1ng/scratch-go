package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var multiplyCmd = &cobra.Command{
	Use:   "multiply num1 num2 [numN]",
	Short: "Multiply subcommand add all passed args",

	Run: func(cmd *cobra.Command, args []string) {
		nums := convertArgsToFloatSlice(args, ErrorHandling(parseHandling))
		result := calc(nums, MULTIPLY)
		fmt.Fprintf(os.Stdout, "%s = %.2f\n", strings.Join(args, "*"), result)
	},
}

func init() {
	rootCmd.AddCommand(multiplyCmd)
}
