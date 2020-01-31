package ad

import (
	"fmt"

	"github.com/dlampsi/ldapconn"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ldap.v3"
)

// Groups struct.
type Groups struct {
	c *Catalog
}

// Groups Entry point to AD catalog groups methods.
func (c *Catalog) Groups() *Groups {
	return &Groups{c: c}
}

// Group Grouper group entry struct.
type Group struct {
	ID            string     `json:"id"`
	DN            string     `json:"dn"`
	CN            string     `json:"cn"`
	Description   string     `json:"description"`
	Members       Membership `json:"members"`
	MembersGroups []string
}

// GetByAccountNameShort Get user info from catalog by user sAMAccountName.
// Returns nil structure if user not found.
func (it *Groups) GetByFilter(filter string) ([]Group, error) {
	log.WithFields(log.Fields{"filter": filter, "base": it.c.base}).Debug("GetByFilter")
	sr := ldapconn.CreateRequest(it.c.base, filter)
	entries, err := it.c.cl.SearchEntries(sr)
	if err != nil {
		return nil, fmt.Errorf("can't perform search: %s", err.Error())
	}
	if len(entries) == 0 {
		return nil, nil
	}
	groups := []Group{}
	for _, e := range entries {
		group := Group{
			DN:          e.DN,
			ID:          e.GetAttributeValue("sAMAccountName"),
			CN:          e.GetAttributeValue("cn"),
			Description: e.GetAttributeValue("description"),
			Members:     Membership{},
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (it *Groups) GetMembers(group *Group, nested bool) (*Group, error) {
	membersFilter := fmt.Sprintf(`(&(objectCategory=person)(memberOf=%v))`, ldap.EscapeFilter(group.DN))
	if nested {
		membersFilter = fmt.Sprintf(`(&(objectCategory=person)(memberOf:1.2.840.113556.1.4.1941:=%v))`, ldap.EscapeFilter(group.DN))
	}
	members, err := it.c.Users().GetByFilter(membersFilter)
	if err != nil {
		return nil, fmt.Errorf("can't get group members: %s", err.Error())
	}
	for _, member := range members {
		group.Members.All = append(group.Members.All, member.ID)
		group.Members.DirectDN = append(group.Members.DirectDN, member.DN)
	}
	group.Members.Count = len(group.Members.All)
	return group, nil
}
