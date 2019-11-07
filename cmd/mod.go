package cmd

import (
	"cataloger/catalogs/ad"
	"fmt"

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
			default:
				log.Errorf("Unknown source '%s'", source)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(modCmd)
	modCmd.AddCommand(modGroupCmd)

	modGroupCmd.Flags().StringSliceVarP(&addMembers, "add-member", "a", []string{}, "Add new members to group")
	modGroupCmd.Flags().StringSliceVarP(&delMembers, "del-member", "d", []string{}, "Delete members from group")
}

// Modify AD group.
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
				log.Warning("No new members to add")
			default:
				log.Fatal(err)
			}
		}
	}
	if len(delMembers) > 0 {
		if err := c.Groups().DelMembersByAccountName(args[0], addMembers); err != nil {
			log.Fatal(err)
		}
	}
}
