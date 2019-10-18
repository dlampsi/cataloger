package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	delCmd = &cobra.Command{
		Use:   "del",
		Short: "Delete data from catalog entries",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}
)

func init() {
	rootCmd.AddCommand(delCmd)
}
