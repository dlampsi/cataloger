package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize client. Login to catalog",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}
