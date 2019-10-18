package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add data to catalog entries",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}
