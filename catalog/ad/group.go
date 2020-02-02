package ad

import (
	"cataloger/client"
	"errors"
	"fmt"
	"sync"

	"github.com/dlampsi/generigo"
	"github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNoNewMembersToAdd = errors.New("no new members to add")
	ErrNoNewMembersToDel = errors.New("no new members to delete")
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
	log.WithFields(log.Fields{"filter": filter, "base": it.c.Attributes.SearchBase}).Debug("GetByFilter")
	sr := client.NewSearchRequest(it.c.Attributes.SearchBase, filter)
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
	filter := fmt.Sprintf(`(&(objectCategory=person)(memberOf=%v))`, ldap.EscapeFilter(group.DN))
	if nested {
		filter = fmt.Sprintf(`(&(objectCategory=person)(memberOf:1.2.840.113556.1.4.1941:=%v))`, ldap.EscapeFilter(group.DN))
	}
	log.WithFields(log.Fields{"filter": filter, "base": it.c.Attributes.SearchBase}).Debug("GetMembers")
	members, err := it.c.Users().GetByFilter(filter)
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

func (it *Groups) AddMembers(g *Group, members []string) error {
	wg := &sync.WaitGroup{}
	ch := make(chan string, 500)
	for _, m := range members {
		wg.Add(1)
		go func(g *Group, m string, wg *sync.WaitGroup, ch chan<- string) {
			defer wg.Done()
			log.Debugf("Processing member '%s'", m)
			users, err := it.c.Users().GetByFilter(fmt.Sprintf("(sAMAccountName=%s)", m))
			if err != nil {
				log.Warnf("Can't get member '%s': %s", m, err.Error())
				return
			}
			if len(users) == 0 {
				log.Warnf("Member entry '%s' not found", m)
				return
			}
			u := users[0]
			if generigo.StringInSlice(u.ID, g.Members.All) {
				log.Warnf("User '%s' already member of a group", m)
				return
			}
			log.Debugf("New member '%s' will be added", m)
			ch <- u.DN
		}(g, m, wg, ch)
	}
	wg.Wait()
	close(ch)
	toadd := make([]string, 0, len(ch))
	for r := range ch {
		toadd = append(toadd, r)
	}
	if len(toadd) == 0 {
		return ErrNoNewMembersToAdd
	}
	log.Debugf("Members to add: %v", toadd)
	newMembersDN := make([]string, 0, len(g.Members.DirectDN)+len(toadd))
	newMembersDN = append(newMembersDN, g.Members.DirectDN...)
	newMembersDN = append(newMembersDN, toadd...)

	if err := it.c.cl.UpdateAttribute(g.DN, "member", newMembersDN); err != nil {
		return fmt.Errorf("error update group: %s", err.Error())
	}
	log.Infof("Added %d group members", len(toadd))
	return nil
}

func (it *Groups) DelMembers(g *Group, members []string) error {
	wg := &sync.WaitGroup{}
	ch := make(chan string, 500)
	for _, m := range members {
		wg.Add(1)
		go func(g *Group, m string, wg *sync.WaitGroup, ch chan<- string) {
			defer wg.Done()
			log.Debugf("Processing member '%s'", m)
			users, err := it.c.Users().GetByFilter(fmt.Sprintf("(sAMAccountName=%s)", m))
			if err != nil {
				log.Warnf("Can't get member '%s': %s", m, err.Error())
				return
			}
			if len(users) == 0 {
				log.Warnf("Member entry '%s' not found", m)
				return
			}
			u := users[0]
			if !generigo.StringInSlice(u.ID, g.Members.All) {
				log.Warnf("User '%s' already not a member of a group", m)
				return
			}
			log.Debugf("Member '%s' will be deleted from group", m)
			ch <- u.DN
		}(g, m, wg, ch)
	}
	wg.Wait()
	close(ch)
	todell := make([]string, 0, len(ch))
	for r := range ch {
		todell = append(todell, r)
	}
	if len(todell) == 0 {
		return ErrNoNewMembersToDel
	}
	log.Debugf("Members to delete: %v", todell)

	// Delete from group only current group members
	newMembersDN := make([]string, 0, len(g.Members.DirectDN))
	for _, curMember := range g.Members.DirectDN {
		if !generigo.StringInSlice(curMember, todell) {
			newMembersDN = append(newMembersDN, curMember)
		}
	}
	log.Debug(newMembersDN)
	if err := it.c.cl.UpdateAttribute(g.DN, "member", newMembersDN); err != nil {
		return fmt.Errorf("error update group: %s", err.Error())
	}
	log.Infof("Removed %d group members", len(todell))
	return nil
}
