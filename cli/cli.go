package neuroncli

import (
	"fmt"
	command "neuron/cli/commands"
        "github.com/spf13/cobra"
	"os"
)

var (
	cmd *cobra.Command
)

func init() {
	cmd = command.SetNeuronCmds()
}

func CliMain() {
	err := Execute(os.Args[1:])
	if err != nil {
		fmt.Println("An error occured")
		os.Exit(1)
	}
}

func Execute(args []string) error {

	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	if err != nil {
		return err
	}
	return nil
}
