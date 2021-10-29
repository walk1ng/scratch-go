package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add num1 num2 [numN]",
	Short: "Add subcommand add all passed args",

	Run: func(cmd *cobra.Command, args []string) {
		nums := convertArgsToFloatSlice(args, ErrorHandling(parseHandling))
		result := calc(nums, ADD)
		fmt.Fprintf(os.Stdout, "%s = %.2f\n", strings.Join(ConvertValuesToStringSlice(nums), "+"), result)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
