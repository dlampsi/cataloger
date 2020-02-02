package cmd

import (
	"cataloger/catalog/ad"
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
	}
	modifyGroupMembersCmd = &cobra.Command{
		Use:   "members [groups_ids]",
		Short: "Modify group members",
		Run:   modifyGroupMembersRun,
	}
)

func init() {
	rootCmd.AddCommand(modifyCmd)
	bindPersistentFlag(modifyCmd, "String", &flagAttributes{
		Id:           "attribute",
		Description:  "Search attribute name for modified entry",
		DefaultValue: "",
	})

	modifyCmd.AddCommand(modifyGroupCmd)

	modifyGroupCmd.AddCommand(modifyGroupMembersCmd)
	bindFlag(modifyGroupMembersCmd, "StringSliceP", &flagAttributes{
		Id:           "add",
		Short:        "a",
		Description:  "List of members to add to group",
		DefaultValue: []string{},
	})
	bindFlag(modifyGroupMembersCmd, "StringSliceP", &flagAttributes{
		Id:           "delete",
		Short:        "d",
		Description:  "List of members to delete from group",
		DefaultValue: []string{},
	})
}

func modifyGroupMembersRun(cmd *cobra.Command, args []string) {
	catalog, err := initAdCatalog()
	if err != nil {
		log.Fatal(err)
	}
	addMembers := viper.GetStringSlice("add")
	delMembers := viper.GetStringSlice("delete")
	if len(addMembers) == 0 && len(delMembers) == 0 {
		log.Warning("Empty add or remove members list")
		os.Exit(2)
	}

	for _, arg := range args {
		filter := fmt.Sprintf("(&(objectClass=group)(%s=%s))", catalog.Attributes.SearchAttribute, arg)
		groups, err := catalog.Groups().GetByFilter(filter)
		if err != nil {
			log.Error(err)
			continue
		}
		if len(groups) == 0 {
			log.Error("No group entries found")
			continue
		}
		if len(groups) == 0 {
			log.Info("No groups provided in command args")
		}
		for _, g := range groups {
			g, err := catalog.Groups().GetMembers(&g, false)
			if err != nil {
				log.Error(err)
				continue
			}
			// Add members
			if len(addMembers) > 0 {
				if err := catalog.Groups().AddMembers(g, addMembers); err != nil {
					if err == ad.ErrNoNewMembersToAdd {
						log.Warn(err)
					} else {
						log.Error(err)
					}
				}
			}
			// Delete members
			if len(delMembers) > 0 {
				if err := catalog.Groups().DelMembers(g, delMembers); err != nil {
					if err == ad.ErrNoNewMembersToDel {
						log.Warn(err)
					} else {
						log.Error(err)
					}
				}
			}
		}
	}
}
