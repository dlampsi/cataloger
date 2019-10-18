package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:   "get [enties_id]",
		Short: "Get data from catalog",
		Run: func(cmd *cobra.Command, args []string) {
			for _, e := range args {
				fmt.Printf("Get entrie info: %s\n", e)
			}
		},
	}

	getUserCmd = &cobra.Command{
		Use:   "users [users_id]",
		Short: "Get users data from catalog",
		Run: func(cmd *cobra.Command, args []string) {
			for _, u := range args {
				fmt.Printf("Get user info: %s\n", u)
			}
		},
	}

	getGroupCmd = &cobra.Command{
		Use:   "groups [groups_id]",
		Short: "Get groups data from catalog",
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
