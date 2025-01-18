package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

var (
	defaultMarkDownFormat = heredoc.Doc(`
        ​---
        title: "%s"
        date: %s
        description: "%s"
        ​---
    `)
	title, description string

	newCmd = &cobra.Command{
		Use:   "new <filename>",
		Short: "Create a markdown file based on a specific format.",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]

			file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
			if err == os.ErrExist {
				return err
			} else if err != nil {
				return err
			}

			defer file.Close()

			if title == "" {
				fmt.Print("title: ")
				fmt.Scanln(&title)
			}

			if description == "" {
				fmt.Print("description: ")
				fmt.Scanln(&description)
			}

			now := time.Now().Format(time.RFC3339)

			if _, err := file.WriteString(fmt.Sprintf(defaultMarkDownFormat, title, now, description)); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	newCmd.Flags().StringVarP(&title, "title", "t", "", "title (if not provided, qq will prompt the user for it)")
	newCmd.Flags().StringVarP(&description, "description", "d", "", "description (if not provided, qq will prompt the user for it)")
}
