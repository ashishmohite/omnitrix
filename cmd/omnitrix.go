package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var omnitrix = &cobra.Command{
	Use:   "omnitrix",
	Short: "Omnitrix is a transformer of DNA🧬 samples to Aliens👽",
}

func Execute() {
	if err := omnitrix.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
