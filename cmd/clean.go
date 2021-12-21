package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var cleanDryRun bool

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean files and folders that are rotten",
	Long:  `Clean files and folders that are rotten.`,
	Run:   clean,
}

var clean = func(cmd *cobra.Command, args []string) {
	var (
		err      error
		deleted  bool
		rotItems []RotItem
	)
	rotItems, err = load()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var rotItemsNew []RotItem
	now := time.Now()
	for _, rotItem := range rotItems {
		deleted, err = rotItem.Clean(now, cleanDryRun)
		if err != nil {
			fmt.Println(err.Error())
		}
		if cleanDryRun || (!deleted || err != nil) {
			rotItemsNew = append(rotItemsNew, rotItem)
		}
		if deleted {
			fmt.Printf("Removed %s.\n", rotItem.Path)
		}
	}

	err = save(rotItemsNew)
	if err != nil {
		fmt.Println("Saving failed: " + err.Error())
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(cleanCmd)

	cleanCmd.Flags().BoolVarP(&cleanDryRun, "dry-run", "d", false, "Only simulate actions")
}
