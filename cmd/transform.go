package cmd

import (
	"fmt"
	"os"

	"omnitrix/dna"

	"github.com/spf13/cobra"
)

var SamplePath string

var transform = &cobra.Command{
	Use: "transform [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
			}
			_ = (&dna.Sample{Config: map[string]interface{}{}, Path: args[0]}).Transform(pwd)
		} else if len(args) == 2 {
			_ = (&dna.Sample{Config: map[string]interface{}{}, Path: args[0]}).Transform(args[1])
		}
	},
}

func init() {
	omnitrix.AddCommand(transform)
	transform.Flags().StringVarP(&SamplePath, "sample", "s", "", "Path to DNA sample(template)")
}
