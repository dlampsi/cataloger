package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cataloger",
		Short: "Util for interact with ldap and active directory catalogs.",
	}

	source string
)

// Execute adds all child commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(onInit)

	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "ad", "Source catalog type. Avalible: ad, ldap.")
}

func onInit() {
}
