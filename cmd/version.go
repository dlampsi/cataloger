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
			if fullPrint {
				fmt.Println(info.ForPrintFull())
				return
			}
			fmt.Println(info.ForPrint())
		},
	}
	fullPrint bool
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().BoolVarP(&fullPrint, "full", "f", false, "Print version and additional information")
}
