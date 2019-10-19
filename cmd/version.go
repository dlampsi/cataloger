package cmd

import (
	"cataloger/info"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints cataloger version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(info.ForPrint())
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
