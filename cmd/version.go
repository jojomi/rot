package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd prints the program version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print current version",
	Run:   versionHandler,
}

var versionHandler = func(cmd *cobra.Command, args []string) {
	fmt.Printf("%v, commit %v, built at %v\n", version, commit, date)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
