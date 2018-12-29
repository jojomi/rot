package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listChangedOnly bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the files and folders current rotting",
	Run:   list,
}

var list = func(cmd *cobra.Command, args []string) {
	rotItems, err := load()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	/*
	 * broken for umlauts:
	 * The Writer assumes that all Unicode code points have the same width; this may not be true in some fonts or if the string contains combining characters.
	 * source: https://golang.org/pkg/text/tabwriter/
	 */
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 1, 3, ' ', 0)

	fmt.Fprintln(w, "ID\tDelete after\tChanged\tDelete if changed\tType\tPath\tCreated rot dataset")
	var (
		dataType, changedString, deleteIfChanged string
		changed, checkedChange                   bool
	)

	for i, rotItem := range rotItems {
		changed = false
		checkedChange = false
		if listChangedOnly {
			changed = rotItem.HasChanged()
			checkedChange = true
			if !changed {
				continue
			}
		}
		if !checkedChange {
			changed = rotItem.HasChanged()
		}

		if rotItem.IsFolder {
			dataType = "Folder"
		} else {
			dataType = "File"
		}
		if changed {
			changedString = "*"
		} else {
			changedString = "-"
		}
		if rotItem.DeleteIfModified {
			deleteIfChanged = "✓"
		} else {
			deleteIfChanged = "×"
		}

		fmt.Fprintln(w,
			fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\t%s",
				i+1,
				rotItem.DeleteAt.Format("2006-01-02 15:05:06"),
				changedString,
				deleteIfChanged,
				dataType,
				rotItem.Path,
				rotItem.AddedAt.Format("2006-01-02 15:05:06"),
			),
		)
	}

	w.Flush()
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&listChangedOnly, "changed", "c", false, "Only list changed files")
}
