package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	modCmd = &cobra.Command{
		Use:   "mod",
		Short: "Modifying operations",
	}

	modGroupCmd = &cobra.Command{
		Use:   "group <group_id>",
		Short: "Group modify operations",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Unexpected number of <group_id> argument; want: 1, get: %d", len(args))
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%v\n", args)
		},
	}
)

func init() {
	rootCmd.AddCommand(modCmd)
	modCmd.AddCommand(modGroupCmd)

	modGroupCmd.Flags().StringSliceP("add-member", "a", []string{}, "Add new members to group")
	modGroupCmd.Flags().StringSliceP("del-member", "d", []string{}, "Delete members from group")
}
