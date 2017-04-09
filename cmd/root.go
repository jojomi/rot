package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// AppFs is a layer of abstraction for the filesystem
var AppFs afero.Fs = afero.NewOsFs()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "rot",
	Short: "rot empowers you to stage files and folders for rotting (later deletion).",
	Long:  `rot empowers you to stage files and folders for rotting (later deletion).`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
}
