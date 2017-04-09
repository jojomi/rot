package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var addDeleteIn string
var addDeleteAt string
var addDeleteIfModified bool

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add files or folders for rotting",
	Long:  `This command adds a file or a folder for later deletion (rotting).`,
	Run:   add,
}

var add = func(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("No file given")
		os.Exit(1)
	}

	var (
		rotItem RotItem
		err     error
	)
	rotItem, err = NewRotItem(args[0], addDeleteIfModified)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if addDeleteIn != "" {
		rotItem.SetDeletionDuration(addDeleteIn)
	}
	if addDeleteAt != "" {
		var date time.Time
		date, _ = time.Parse("2006-01-02 15:05:06", addDeleteAt)
		if date.IsZero() {
			date, _ = time.Parse("2006-01-02", addDeleteAt+" 23:59:59")
		}
		if !date.IsZero() {
			rotItem.SetDeletionDate(date)
		}
	}
	if rotItem.DeleteAt.IsZero() {
		fmt.Println("No valid date given for deletion (use either --at or --in)")
		os.Exit(2)
	}

	data, err := load()
	if err != nil {
		fmt.Println("loading failed")
		os.Exit(1)
	}
	data = append(data, rotItem)
	err = save(data)
	if err != nil {
		fmt.Println("Saving failed")
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addDeleteIn, "in", "i", "", "Delete file or folder in a certain time from now. Format as ISO 8601, e.g. P4D to delete in 4 days.")
	addCmd.Flags().StringVarP(&addDeleteAt, "at", "a", "", "Delete file or folder at a certain time. Format like 2006-01-02 15:05:06.")
	addCmd.Flags().BoolVarP(&addDeleteIfModified, "delete-if-modified", "m", false, "Delete file or folder even if it was modified meanwhile (default: false).")
}
