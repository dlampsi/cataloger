package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove catalog entries",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}
)

func init() {
	rootCmd.AddCommand(removeCmd)
}
