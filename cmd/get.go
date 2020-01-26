package cmd

import (
	"cataloger/catalogs/ad"
	"cataloger/catalogs/ldap"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	// Default attribute for get entries AD
	defaultSearchAttrAD = "sAMAccountName"
	// Default attribute for get entries LDAP
	defaultSearchAttrLDAP = "uid"
)

var (
	getShort              bool
	showOnlyGroups        bool
	showOnlyGroupsDirect  bool
	showOnlyMembers       bool
	showOnlyMembersDirect bool
	searchAttr            string

	getCmd = &cobra.Command{
		Use:   "get [enties_id]",
		Short: "Get data from catalog",
	}

	getUserCmd = &cobra.Command{
		Use:   "user [users_id]",
		Short: "Get users data from catalog",
		Run: func(cmd *cobra.Command, args []string) {
			switch source {
			case "ad":
				getUsersAd(args)
			case "ldap":
				getUsersLdap(args)
			default:
				log.Errorf("Unknown source '%s'", source)
				os.Exit(1)
			}
		},
	}

	getGroupCmd = &cobra.Command{
		Use:   "group [groups_id]",
		Short: "Get group data from catalog",
		Run: func(cmd *cobra.Command, args []string) {
			switch source {
			case "ad":
				getGroupsAd(args)
			case "ldap":
				getGroupsLdap(args)
			default:
				log.Errorf("Unknown source '%s'", source)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().BoolVar(&getShort, "short", false, "Print short user info")
	getCmd.PersistentFlags().StringVarP(&searchAttr, "search-attribute", "a", "",
		"Specifies search attribute. Defaults for AD - 'sAMAccountName', for LDAP - 'uid'")

	getCmd.AddCommand(getUserCmd)
	getUserCmd.Flags().BoolVarP(&showOnlyGroups, "groups", "g", false, "Print only user groups")
	getUserCmd.Flags().BoolVarP(&showOnlyGroupsDirect, "direct", "d", false, "Print only direct user groups (only AD)")

	getCmd.AddCommand(getGroupCmd)
	getGroupCmd.Flags().BoolVarP(&showOnlyMembers, "members", "m", false, "Print only group members")
	getGroupCmd.Flags().BoolVarP(&showOnlyMembersDirect, "direct", "d", false, "Print only direct group members")
}

func getUsersAd(args []string) {
	c, err := ad.NewCatalog(createConfig())
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range args {
		log.Debugf("Get user info: %s\n", u)
		var data *ad.UserEntry
		if getShort {
			data, err = c.Users().GetByAccountNameShort(u)
		} else {
			data, err = c.Users().GetByAccountName(u)
		}

		if err != nil {
			log.Errorf("ERROR: %s\n", err.Error())
		}
		if data == nil {
			log.Errorf("Entry '%s' not found\n", u)
		} else {
			if showOnlyGroups {
				c.Users().PrintGroups(data)
			} else if showOnlyGroupsDirect {
				c.Users().PrintGroupsDirect(data)
			} else {
				c.Users().Print(data)
			}
		}
	}
}

func getGroupsAd(args []string) {
	c, err := ad.NewCatalog(createConfig())
	if err != nil {
		log.Fatal(err)
	}

	for _, g := range args {
		log.Debugf("Get group info: %s\n", g)
		var data *ad.GroupEntry
		if getShort {
			data, err = c.Groups().GetByAccountNameShort(g)
		} else {
			data, err = c.Groups().GetByAccountName(g)
		}
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
		if data == nil {
			fmt.Printf("Entry '%s' not found\n", g)
		} else {
			if showOnlyMembers {
				c.Groups().PrintMembers(data)
			} else if showOnlyMembersDirect {
				c.Groups().PrintMembersDirect(data)
			} else {
				c.Groups().Print(data)
			}
		}
	}
}

func getUsersLdap(args []string) {
	c, err := ldap.NewCatalog(createConfig())
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range args {
		log.Debugf("Get user info: %s\n", u)
		var data *ldap.UserEntry
		if getShort {
			data, err = c.Users().GetByCnShort(u)
		} else {
			data, err = c.Users().GetByCn(u)
		}

		if err != nil {
			log.Errorf("ERROR: %s\n", err.Error())
		}
		if data == nil {
			log.Errorf("Entry '%s' not found\n", u)
		} else {
			if showOnlyGroups {
				c.Users().PrintGroups(data)
			} else {
				c.Users().Print(data)
			}
		}
	}
}

func getGroupsLdap(args []string) {
	for _, g := range args {
		log.Debugf("Get group info: %s\n", g)
		c, err := ldap.NewCatalog(createConfig())
		if err != nil {
			log.Fatal(err)
		}
		data, err := c.Groups().GetByCn(g)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
		if data == nil {
			fmt.Printf("Entry '%s' not found\n", g)
			continue
		}
		if showOnlyMembers {
			c.Groups().PrintMembers(data)
		} else {
			c.Groups().Print(data)
		}
	}
}
