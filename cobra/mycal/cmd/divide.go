package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	dividedByZeroHandling int
)

var divideCmd = &cobra.Command{
	Use:   "divide num1 num2 [numN]",
	Short: "Divide subcommand add all passed args",

	Run: func(cmd *cobra.Command, args []string) {
		nums := convertArgsToFloatSlice(args, ErrorHandling(parseHandling))
		result := calc(nums, DIVIDE)
		fmt.Fprintf(os.Stdout, "%s = %.2f\n", strings.Join(ConvertValuesToStringSlice(nums), "/"), result)
	},
}

func init() {

	divideCmd.Flags().IntVarP(&dividedByZeroHandling, "divide_by_zero", "d", int(PanicOnDividedByZero), "what behavior when divided by zero error")

	rootCmd.AddCommand(divideCmd)
}
