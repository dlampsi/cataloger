package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove entries from catalog",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%v\n", args)
		},
	}
)

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.PersistentFlags().String("id", "", "Entry ID attribute")
	rmCmd.PersistentFlags().String("dn", "", "Entry DN attribute")
}
