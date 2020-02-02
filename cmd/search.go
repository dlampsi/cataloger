package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	searchCmd = &cobra.Command{
		Use:   "search",
		Short: "Search for entries",
	}
	searchUserCmd = &cobra.Command{
		Use:   "user [users_ids]",
		Short: "Search for user entires",
		Run:   searchUserRun,
	}
	searchGroupCmd = &cobra.Command{
		Use:   "group [groups_ids]",
		Short: "Search for group entires",
		Run:   searchGroupRun,
	}
)

func init() {
	rootCmd.AddCommand(searchCmd)
	bindPersistentFlag(searchCmd, "String", &flagAttributes{
		Id:           "filter",
		Description:  "LDAP search filter",
		DefaultValue: "",
	})
	bindPersistentFlag(searchCmd, "String", &flagAttributes{
		Id:           "attribute",
		Description:  "Search attribute name",
		DefaultValue: "",
	})

	searchCmd.AddCommand(searchUserCmd)
	bindFlag(searchUserCmd, "BoolP", &flagAttributes{
		Id:          "user-groups",
		Short:       "g",
		Description: "Search user groups memberships",
	})

	searchCmd.AddCommand(searchGroupCmd)
	bindFlag(searchGroupCmd, "BoolP", &flagAttributes{
		Id:          "group-members",
		Short:       "m",
		Description: "Search group members",
	})
	bindFlag(searchGroupCmd, "BoolP", &flagAttributes{
		Id:          " ",
		Short:       "n",
		Description: "Search for nested grop members",
	})
}

func searchUserRun(cmd *cobra.Command, args []string) {
	catalog, err := initAdCatalog()
	if err != nil {
		log.Fatal(err)
	}
	filter := viper.GetString("filter")
	for _, arg := range args {
		if filter == "" {
			filter = fmt.Sprintf("(&(objectClass=person)(%s=%s))", catalog.Attributes.SearchAttribute, arg)
		}
		users, err := catalog.Users().GetByFilter(filter)
		if err != nil {
			log.Error(err)
			continue
		}
		if users == nil {
			log.Error("No users found")
			continue
		}
		for _, u := range users {
			if viper.GetBool("user-groups") {
				ug, err := catalog.Users().GetGroups(&u)
				if err != nil {
					log.Error(err)
					continue
				}
				catalog.Users().Printer(ug)
			} else {
				catalog.Users().Printer(&u)
			}
		}
	}
}

func searchGroupRun(cmd *cobra.Command, args []string) {
	catalog, err := initAdCatalog()
	if err != nil {
		log.Fatal(err)
	}
	filter := viper.GetString("filter")
	for _, arg := range args {
		if filter == "" {
			filter = fmt.Sprintf("(&(objectClass=group)(%s=%s))", catalog.Attributes.SearchAttribute, arg)
		}
		groups, err := catalog.Groups().GetByFilter(filter)
		if err != nil {
			log.Error(err)
			continue
		}
		if groups == nil {
			log.Error("No groups found")
			continue
		}
		for _, g := range groups {
			g, err := catalog.Groups().GetMembers(&g, viper.GetBool("nested"))
			if err != nil {
				log.Error(err)
			} else {
				catalog.Groups().Printer(g, viper.GetBool("group-members"))
			}
		}
	}
}
