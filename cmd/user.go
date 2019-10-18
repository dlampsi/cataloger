package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "User entrie operations",
	}

	userCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create new user entrie",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}

	userDelCmd = &cobra.Command{
		Use:   "del",
		Short: "Delete user entrie",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}

	userAddToCmd = &cobra.Command{
		Use:   "add-to [groups_ids]",
		Short: "Add user to group members",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}

	userDelFromCmd = &cobra.Command{
		Use:   "del-from [groups_ids]",
		Short: "Delete user from a groups members",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}
)

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userDelCmd)
	userCmd.AddCommand(userAddToCmd)
	userCmd.AddCommand(userDelFromCmd)
}
