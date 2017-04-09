package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops a file or folder from rotting",
	Long:  `Stops a file or folder from rotting.`,
	Run:   stop,
}

var stop = func(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("No ID given.")
		os.Exit(1)
	}

	ID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("No valid ID given.")
		os.Exit(1)
	}
	data, err := load()
	if err != nil {
		fmt.Println("Load failed.")
		os.Exit(1)
	}
	data = append(data[:ID-1], data[ID:]...)
	err = save(data)
	if err != nil {
		fmt.Println("Saving failed.")
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
