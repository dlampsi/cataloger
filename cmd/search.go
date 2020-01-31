package cmd

import (
	"cataloger/ad"
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
		Use:   "user",
		Short: "Search for user entires",
		Run: func(cmd *cobra.Command, args []string) {
			adSearchUser(args)
		},
	}
	searchGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Search for group entires",
		Run: func(cmd *cobra.Command, args []string) {
			adSearchGroup(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.PersistentFlags().String("filter", "", "LDAP search filter")
	bindSearchFlag("filter")
	searchCmd.PersistentFlags().StringP("attribute", "a", "", "Search attribute name")
	bindSearchFlag("attribute")

	searchCmd.AddCommand(searchUserCmd)
	searchUserCmd.Flags().BoolP("user-groups", "g", false, "Search user groups memberships")
	if err := viper.BindPFlag("user-groups", searchUserCmd.Flags().Lookup("user-groups")); err != nil {
		log.Fatal(err)
	}

	searchCmd.AddCommand(searchGroupCmd)
	searchGroupCmd.Flags().BoolP("group-members", "m", false, "Search group members")
	if err := viper.BindPFlag("group-members", searchGroupCmd.Flags().Lookup("group-members")); err != nil {
		log.Fatal(err)
	}
	searchGroupCmd.Flags().BoolP("nested", "n", false, "Search for nested grop members")
	if err := viper.BindPFlag("nested", searchGroupCmd.Flags().Lookup("nested")); err != nil {
		log.Fatal(err)
	}
}

func bindSearchFlag(flagId string) {
	if err := viper.BindPFlag(flagId, searchCmd.PersistentFlags().Lookup(flagId)); err != nil {
		log.Fatal(err)
	}
}

func adSearchUser(args []string) {
	ldapCfg := createLdapClient()

	c, err := ad.NewCatalog(ldapCfg, viper.GetString("search-base"))
	if err != nil {
		log.Fatalf("Error connect to catalog: %s", err.Error())
	}
	filter := viper.GetString("filter")
	attribute := viper.GetString("attribute")
	if attribute == "" {
		attribute = "sAMAccountName"
	}

	for _, arg := range args {
		if filter == "" {
			filter = fmt.Sprintf("(&(objectClass=person)(%s=%s))", attribute, arg)
		}
		users, err := c.Users().GetByFilter(filter)
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
				ug, err := c.Users().GetGroups(&u)
				if err != nil {
					log.Error(err)
					continue
				}
				c.Users().Printer(ug)
			} else {
				c.Users().Printer(&u)
			}
		}
	}
}

func adSearchGroup(args []string) {
	ldapCfg := createLdapClient()

	c, err := ad.NewCatalog(ldapCfg, viper.GetString("search-base"))
	if err != nil {
		log.Fatalf("Error connect to catalog: %s", err.Error())
	}
	filter := viper.GetString("filter")
	attribute := viper.GetString("attribute")
	if attribute == "" {
		attribute = "sAMAccountName"
	}

	for _, arg := range args {
		if filter == "" {
			filter = fmt.Sprintf("(&(objectClass=group)(%s=%s))", attribute, arg)
		}
		groups, err := c.Groups().GetByFilter(filter)
		if err != nil {
			log.Error(err)
			continue
		}
		if groups == nil {
			log.Error("No groups found")
			continue
		}
		for _, g := range groups {
			g, err := c.Groups().GetMembers(&g, viper.GetBool("nested"))
			if err != nil {
				log.Error(err)
			} else {
				c.Groups().Printer(g, viper.GetBool("group-members"))
			}
		}
	}
}
