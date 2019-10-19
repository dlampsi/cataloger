package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Create new entries in catalog",
	}

	addGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Create new group entry",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Create new group entry...")
		},
	}

	addUserCmd = &cobra.Command{
		Use:   "user",
		Short: "Create new group entry",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Create new group entry...")
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.AddCommand(addUserCmd)
	// TODO: user add flags
	addCmd.AddCommand(addGroupCmd)
	// TODO: group add flags
}
