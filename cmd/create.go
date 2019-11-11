package cmd

import (
	"cataloger/catalogs/ldap"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create new entries in catalog",
	}

	createGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Create new group entry",
		Run: func(cmd *cobra.Command, args []string) {
			switch source {
			case "ad":
				log.Error("Sorry, method not implemented yet")
				os.Exit(2)
			case "ldap":
				createGroupLdap()
			default:
				log.Errorf("Unknown source '%s'", source)
				os.Exit(1)
			}
		},
	}
	groupCn    string
	groupGid   string
	groupDescr string
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.AddCommand(createGroupCmd)
	createGroupCmd.Flags().StringVar(&groupCn, "cn", "", "Group cn")
	createGroupCmd.Flags().StringVar(&groupGid, "gid", "", "Group gid")
	createGroupCmd.Flags().StringVarP(&groupDescr, "description", "d", "Created from cataloger", "Group description")
}

func createGroupLdap() {
	c, err := ldap.NewCatalog(createConfig())
	if err != nil {
		log.Fatal(err)
	}
	g, err := c.Groups().Create(groupCn, groupGid, groupDescr)
	if err != nil {
		switch err {
		case ldap.ErrAlreadyExists:
			log.Warnf("Group '%s' already exists\n", groupCn)
			os.Exit(2)
		default:
			log.Fatal(err)
		}
	}
	c.Groups().Print(g)
}
