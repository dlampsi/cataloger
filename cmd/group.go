package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	groupCmd = &cobra.Command{
		Use:   "group",
		Short: "Group entrie operations",
	}

	groupCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create new group entrie",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}

	groupDelCmd = &cobra.Command{
		Use:   "del",
		Short: "Delete group entrie",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}

	groupAddMemberCmd = &cobra.Command{
		Use:   "add-member",
		Short: "Add new group members",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}

	groupDelMemberCmd = &cobra.Command{
		Use:   "del-member",
		Short: "Delete members from a group",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}
)

func init() {
	rootCmd.AddCommand(groupCmd)
	groupCmd.AddCommand(groupCreateCmd)
	groupCmd.AddCommand(groupDelCmd)
	groupCmd.AddCommand(groupAddMemberCmd)
	groupCmd.AddCommand(groupDelMemberCmd)
}
