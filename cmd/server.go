package cmd

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

var (
	sercverWatch bool
	serverPort   int
	serverCmd    = &cobra.Command{
		Use:   "server",
		Short: "Start a server and serve files from disk",
		Long: heredoc.Doc(`
            Start a server and serve files from disk.
            It includes a renderer for certain file extensions (e.g., markdown files, image files, etc.).
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
)

func init() {
	serverCmd.Flags().BoolVarP(&sercverWatch, "watch", "w", true, "watch filesystem for changes and do live realoding")
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 1234, "port for server listening")
}
