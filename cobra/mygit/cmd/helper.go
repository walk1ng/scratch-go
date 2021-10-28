package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func ExecuteCommand(name string, subname string, args ...string) (string, error) {
	args = append([]string{subname}, args...)

	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func Error(cmd *cobra.Command, args []string, err error) {
	fmt.Fprintf(os.Stdout, "execute %s args: %v error: %v\n", cmd.Name(), args, err)
	os.Exit(-1)
}
