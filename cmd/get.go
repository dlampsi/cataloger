package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:   "get [enties_id]",
		Short: "Get data from catalog",
	}

	getUserCmd = &cobra.Command{
		Use:   "user [users_id]",
		Short: "Get users data from catalog",
		Run: func(cmd *cobra.Command, args []string) {
			for _, u := range args {
				fmt.Printf("Get user info: %s\n", u)
			}
		},
	}

	getGroupCmd = &cobra.Command{
		Use:   "group [groups_id]",
		Short: "Get group data from catalog",
		Run: func(cmd *cobra.Command, args []string) {
			for _, g := range args {
				fmt.Printf("Get group info: %s\n", g)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getUserCmd)
	getCmd.AddCommand(getGroupCmd)
}
