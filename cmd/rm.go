package cmd

import (
	"cataloger/catalogs/ldap"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove entries from catalog",
	}

	rmGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Remove group entry from catalog",
		Run: func(cmd *cobra.Command, args []string) {
			switch source {
			case "ad":
				log.Error("Sorry, method not implemented yet")
				os.Exit(2)
			case "ldap":
				rmGroupsLdap(args)
			default:
				log.Errorf("Unknown source '%s'", source)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.AddCommand(rmGroupCmd)
}

func rmGroupsLdap(args []string) {
	c, err := ldap.NewCatalog(createConfig())
	if err != nil {
		log.Fatal(err)
	}
	for _, g := range args {
		log.Debugf("Delete group: %s\n", g)
		if err := c.Groups().Delete(g); err != nil {
			switch err {
			case ldap.ErrAlreadyNotExists:
				log.Warnf("Group '%s' already not exists", g)
				os.Exit(2)
			default:
				log.Fatal(err)
			}
		}
		log.Infof("Group '%s' deleted", g)
	}
}
