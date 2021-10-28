package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ExecuteCommand(name string, subname string, args ...string) (string, error) {
	args = append([]string{subname}, args...)

	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func Error(cmd *cobra.Command, args []string, err error) {
	fmt.Fprintf(os.Stderr, "execute %s args: %v error: %v\n", cmd.Name(), args, err)
	os.Exit(-1)
}

func calc(values []float64, opType OpType) float64 {
	var result float64
	if len(values) == 0 {
		return result
	}

	result = values[0]
	for i := 1; i < len(values); i++ {
		switch opType {
		case ADD:
			result += values[i]
		case MINUS:
			result -= values[i]
		case MULTIPLY:
			result *= values[i]
		case DIVIDE:
			if values[i] == 0 {
				switch ErrorHandling(dividedByZeroHandling) {
				case ReturnOnDividedByZero:
					return result
				case PanicOnDividedByZero:
					panic(errors.New("divided by zero"))
				}
			}
		}
	}

	return result
}

func convertArgsToFloatSlice(args []string, errHandling ErrorHandling) []float64 {
	var result []float64 = make([]float64, 0, len(args))

	for _, arg := range args {
		value, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			switch errHandling {
			case ExitOnParseError:
				fmt.Fprintf(os.Stderr, "invalid num: %s\n", arg)
				os.Exit(-1)
			case PanicOnParseError:
				panic(err)
			}
		} else {
			result = append(result, value)
		}
	}

	fmt.Println("args:", result)
	return result
}
