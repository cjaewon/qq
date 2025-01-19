package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
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
			// fmt.Println(contentPath)
			// if len(args) > 1 {
			// 	contentPath = args[0]
			// }

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Printf("Request URL: %s\n", r.URL.Path)

				// 응답 전송
				fmt.Fprintf(w, "You accessed %s", r.URL.Path)
			})

			fmt.Println("http server started on :" + strconv.Itoa(serverPort))

			if err := http.ListenAndServe(":"+strconv.Itoa(serverPort), nil); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	serverCmd.Flags().BoolVarP(&serverWatch, "watch", "w", true, "watch filesystem for changes and do live realoding")
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 1234, "port for server listening")
}
