package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create catalog entries",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}
	createUserCmd = &cobra.Command{
		Use:   "user",
		Short: "Create user entrie",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}
	createGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Create group entrie",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(source)
		},
	}
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createUserCmd)
	createCmd.AddCommand(createGroupCmd)
}
