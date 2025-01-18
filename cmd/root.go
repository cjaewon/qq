package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "qq",
		Short: "qq is a simple markdown viewer similar to Github file explore.",
		Long: heredoc.Doc(`
            qq is a simple markdown viewer similar to Github file explore.
            It also provides useful features, including creating a markdown file based on a specific format. 
            Github repo: github.com/cjaewon/qq
        `),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Do you want help? Use -h, --help flags")
		},
	}
)

func init() {
	rootCmd.AddCommand(newCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
