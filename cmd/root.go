package cmd

import (
	"cataloger/info"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cataloger",
		Short: "Util for interact with ldap and active directory catalogs.",
		Run: func(cmd *cobra.Command, args []string) {
			if fullVersion {
				fmt.Println(info.ForPrintFull())
			}
		},
	}

	source      string
	fullVersion bool
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
	rootCmd.Flags().BoolVar(&fullVersion, "version", false, "Prints full cataloger version")
}

func onInit() {
}
