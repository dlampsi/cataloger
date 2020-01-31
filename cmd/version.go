package cmd

import (
	"fmt"

	"cataloger/info"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Cataloger",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(info.ForPrint())
	},
}
