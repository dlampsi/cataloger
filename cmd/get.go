package cmd

import (
	"cataloger/catalogs/ad"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	getShort              bool
	showOnlyGroups        bool
	showOnlyGroupsDirect  bool
	showOnlyMembers       bool
	showOnlyMembersDirect bool

	getCmd = &cobra.Command{
		Use:   "get [enties_id]",
		Short: "Get data from catalog",
	}

	getUserCmd = &cobra.Command{
		Use:   "user [users_id]",
		Short: "Get users data from catalog",
		Run: func(cmd *cobra.Command, args []string) {
			for _, u := range args {
				log.Debugf("Get user info: %s\n", u)
				if source == "ad" {
					c, err := ad.NewCatalog(createConfig())
					if err != nil {
						log.Fatal(err)
					}

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
		},
	}

	getGroupCmd = &cobra.Command{
		Use:   "group [groups_id]",
		Short: "Get group data from catalog",
		Run: func(cmd *cobra.Command, args []string) {
			for _, g := range args {
				log.Debugf("Get group info: %s\n", g)
				if source == "ad" {
					c, err := ad.NewCatalog(createConfig())
					if err != nil {
						log.Fatal(err)
					}
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
		},
	}
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().BoolVar(&getShort, "short", false, "Print short user info")

	getCmd.AddCommand(getUserCmd)
	getUserCmd.Flags().BoolVarP(&showOnlyGroups, "groups", "g", false, "Print only user groups")
	getUserCmd.Flags().BoolVarP(&showOnlyGroupsDirect, "direct", "d", false, "Print only direct user groups")

	getCmd.AddCommand(getGroupCmd)
	getGroupCmd.Flags().BoolVarP(&showOnlyMembers, "members", "m", false, "Print only group members")
	getGroupCmd.Flags().BoolVarP(&showOnlyMembersDirect, "direct", "d", false, "Print only direct group members")
}
