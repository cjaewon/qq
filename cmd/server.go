package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cjaewon/qq/server"
	"github.com/spf13/cobra"
)

var (
	serverWatch bool
	serverPort  int

	serverCmd = &cobra.Command{
		Use:   "server [path]",
		Short: "Start a server and serve files from disk",
		Long: heredoc.Doc(`
            Start a server and serve files from disk.
            It includes a renderer for certain file extensions (e.g., markdown files, image files, etc.).
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			contentRootPath := "."

			if len(args) >= 1 {
				contentRootPath = args[0]
			}

			f, err := os.Stat(contentRootPath)
			if err != nil {
				return err
			}

			if !f.IsDir() {
				return errors.New(fmt.Sprintf("Error: %s is not directory", contentRootPath))
			}

			s := server.Server{
				Port:               serverPort,
				Watch:              serverWatch,
				ContentRootPath:    contentRootPath,
				ContentRootDirName: f.Name(),
			}

			return s.Start()
		},
	}
)

func init() {
	serverCmd.Flags().BoolVarP(&serverWatch, "watch", "w", true, "watch filesystem for changes and do live realoding")
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 1234, "port for server listening")
}
