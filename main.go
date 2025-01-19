package main

import (
	"os"

	"github.com/cjaewon/qq/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
