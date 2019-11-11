package cmd

import (
	"cataloger/catalogs/ad"
	"cataloger/catalogs/ldap"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	addMembers []string
	delMembers []string

	modCmd = &cobra.Command{
		Use:   "mod",
		Short: "Modifying operations",
	}

	modGroupCmd = &cobra.Command{
		Use:   "group <group_id>",
		Short: "Group modify operations",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Unexpected number of <group_id> argument; want: 1, get: %d", len(args))
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Debugf("Source '%s'", source)
			log.Debugf("Modifying group '%s'", args[0])
			switch source {
			case "ad":
				modGroupAd(args)
			case "ldap":
				modGroupLdap(args)
			default:
				log.Errorf("Unknown source '%s'", source)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(modCmd)
	modCmd.AddCommand(modGroupCmd)

	modGroupCmd.Flags().StringSliceVarP(&addMembers, "add-member", "a", []string{}, "Members ID to add to group")
	modGroupCmd.Flags().StringSliceVarP(&delMembers, "del-member", "d", []string{}, "Members ID to remove from group")
}

func modGroupAd(args []string) {
	c, err := ad.NewCatalog(createConfig())
	if err != nil {
		log.Fatal(err)
	}
	if len(addMembers) > 0 {
		if err := c.Groups().AddMembersByAccountName(args[0], addMembers); err != nil {
			switch err {
			case ad.ErrEntryNotFound:
				log.Fatal("Group not found")
			case ad.ErrEmptyMembersList:
				log.Fatal("Empty add members list (--add-member flags)")
			case ad.ErrNoNewMembersToAdd:
				log.Warning("No new group members to add")
				os.Exit(2)
			default:
				log.Fatal(err)
			}
		}
	}
	if len(delMembers) > 0 {
		if err := c.Groups().DelMembersByAccountName(args[0], delMembers); err != nil {
			switch err {
			case ad.ErrEntryNotFound:
				log.Fatal("Group not found")
			case ad.ErrEmptyMembersList:
				log.Fatal("Empty del members list (--del-member flags)")
			case ad.ErrNoNewMembersToDel:
				log.Warning("No new group members to delete")
				os.Exit(2)
			default:
				log.Fatal(err)
			}
		}
	}
}

func modGroupLdap(args []string) {
	c, err := ldap.NewCatalog(createConfig())
	if err != nil {
		log.Fatal(err)
	}
	if len(addMembers) > 0 {
		if _, err := c.Groups().AddMembers(args[0], addMembers); err != nil {
			switch err {
			case ldap.ErrEntryNotFound:
				log.Fatal("Group not found")
			case ldap.ErrEmptyMembersList:
				log.Fatal("Empty add members list (--add-member flags)")
			case ldap.ErrNoNewMembersToAdd:
				log.Warning("No new group members to add")
				os.Exit(2)
			default:
				log.Fatal(err)
			}
		}
	}
	if len(delMembers) > 0 {
		if _, err := c.Groups().DelMembers(args[0], delMembers); err != nil {
			switch err {
			case ldap.ErrEntryNotFound:
				log.Fatal("Group not found")
			case ldap.ErrEmptyMembersList:
				log.Fatal("Empty del members list (--del-member flags)")
			case ldap.ErrNoNewMembersToDel:
				log.Warning("No new group members to delete")
				os.Exit(2)
			default:
				log.Fatal(err)
			}
		}
	}
}
