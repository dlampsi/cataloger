package cmd

import (
	"cataloger/ad"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	modifyCmd = &cobra.Command{
		Use:   "modify",
		Short: "Search for entries",
	}
	modifyGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Modify group entiry",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	modifyGroupMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "Modify group members",
		Run: func(cmd *cobra.Command, args []string) {
			modGroupAd(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(modifyCmd)
	// modifyCmd.PersistentFlags().String("filter", "", "LDAP search filter")
	// bindModifyFlag("filter")
	modifyCmd.PersistentFlags().String("attribute", "", "Search attribute name for modified entry")
	bindModifyFlag("attribute")

	modifyCmd.AddCommand(modifyGroupCmd)

	modifyGroupCmd.AddCommand(modifyGroupMembersCmd)
	modifyGroupMembersCmd.Flags().StringSliceP("add", "a", []string{}, "List of members to add to group")
	if err := viper.BindPFlag("add", modifyGroupMembersCmd.Flags().Lookup("add")); err != nil {
		log.Fatal(err)
	}
	modifyGroupMembersCmd.Flags().StringSliceP("delete", "d", []string{}, "List of members to delete from group")
	if err := viper.BindPFlag("delete", modifyGroupMembersCmd.Flags().Lookup("delete")); err != nil {
		log.Fatal(err)
	}
}

func bindModifyFlag(flagId string) {
	if err := viper.BindPFlag(flagId, modifyCmd.PersistentFlags().Lookup(flagId)); err != nil {
		log.Fatal(err)
	}
}

func modGroupAd(args []string) {
	ldapCfg := createLdapClient()

	c, err := ad.NewCatalog(ldapCfg, viper.GetString("search-base"))
	if err != nil {
		log.Fatalf("Error connect to catalog: %s", err.Error())
	}
	attribute := viper.GetString("attribute")
	if attribute == "" {
		attribute = "sAMAccountName"
	}

	addMembers := viper.GetStringSlice("add")
	delMembers := viper.GetStringSlice("delete")
	if len(addMembers) == 0 && len(delMembers) == 0 {
		log.Warning("Empty add or remove members list")
		os.Exit(2)
	}

	for _, arg := range args {
		filter := fmt.Sprintf("(&(objectClass=group)(%s=%s))", attribute, arg)
		groups, err := c.Groups().GetByFilter(filter)
		if err != nil {
			log.Error(err)
			continue
		}
		if len(groups) == 0 {
			log.Error("No group entries found")
			continue
		}
		for _, g := range groups {
			g, err := c.Groups().GetMembers(&g, false)
			if err != nil {
				log.Error(err)
				continue
			}
			// Add members
			if len(addMembers) > 0 {
				if err := c.Groups().AddMembers(g, addMembers); err != nil {
					if err == ad.ErrNoNewMembersToAdd {
						log.Warn(err)
					} else {
						log.Error(err)
					}
				}
			}
			// Delete members
			if len(delMembers) > 0 {
				log.Info("delete")
			}
		}
	}
}
