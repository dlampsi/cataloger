package ldap

import (
	"fmt"
	"sync"

	"github.com/dlampsi/generigo"
	"github.com/dlampsi/ldapconn"
	log "github.com/sirupsen/logrus"
)

// Groups struct.
type Groups struct {
	c *Catalog
}

// Groups Entry point to LDAP catalog groups methods.
func (c *Catalog) Groups() *Groups {
	return &Groups{c: c}
}

// GroupEntry Grouper group entry struct.
type GroupEntry struct {
	ID          string          `json:"id"`
	DN          string          `json:"dn"`
	CN          string          `json:"cn"`
	Description string          `json:"description"`
	Members     EntryMembership `json:"members"`
}

// GetByCn Returns LDAP group info.
// Returns nil structure if group not found.
func (it *Groups) GetByCn(cn string) (*GroupEntry, error) {
	filter := fmt.Sprintf("(&(objectClass=posixGroup)(cn=%s))", cn)
	sr := ldapconn.CreateRequest(it.c.groupSearchBase, filter)
	ldapEntry, err := it.c.cl.SearchEntry(sr)
	if err != nil {
		return nil, fmt.Errorf("error get group entry: " + err.Error())
	}
	if ldapEntry == nil {
		return nil, nil
	}

	return &GroupEntry{
		DN:          ldapEntry.DN,
		ID:          ldapEntry.GetAttributeValue("cn"),
		CN:          ldapEntry.GetAttributeValue("cn"),
		Description: ldapEntry.GetAttributeValue("description"),
		Members: EntryMembership{
			Count:    len(ldapEntry.GetAttributeValues("memberUid")),
			Direct:   ldapEntry.GetAttributeValues("memberUid"),
			DirectDN: ldapEntry.GetAttributeValues("memberUid"),
		},
	}, nil
}

// GroupsItems structure.
type GroupsItems struct {
	Base     string   `json:"base"`
	GroupsID []string `json:"groups_id"`
	GroupsDN []string `json:"groups_dn"`
}

// GetAll Returns all LDAP groups in group search base.
func (it *Groups) GetAll() (*GroupsItems, error) {
	filter := "(objectClass=posixGroup)"
	sr := ldapconn.CreateRequest(it.c.searchBase, filter)
	ldapEntries, err := it.c.cl.SearchEntries(sr)
	if err != nil {
		return nil, err
	}

	result := &GroupsItems{
		Base: it.c.searchBase,
	}
	for _, le := range ldapEntries {
		result.GroupsDN = append(result.GroupsDN, le.DN)
		result.GroupsID = append(result.GroupsID, le.GetAttributeValue("cn"))
	}
	return result, nil
}

// Create Creates new LDAP group.
// Returns created ldap group item.
// Returns ErrAlreadyExists if group already exists.
func (it *Groups) Create(cn string, gid string, descr string) (*GroupEntry, error) {
	// Check if group already exists
	g, err := it.GetByCn(cn)
	if err != nil {
		return nil, err
	}
	if g != nil {
		return nil, ErrAlreadyExists
	}

	attr := map[string][]string{
		"objectClass": []string{"top", "posixGroup"},
		"cn":          []string{cn},
		"gidNumber":   []string{gid},
		"description": []string{descr},
	}
	dn := fmt.Sprintf("cn=%s,%s", cn, it.c.groupSearchBase)

	if err := it.c.cl.CreateEntry(dn, attr); err != nil {
		return nil, fmt.Errorf("error create group: %s", err)
	}

	// Get group after create
	created, err := it.GetByCn(cn)
	if err != nil {
		return nil, fmt.Errorf("error get entry after create: %s", err.Error())
	}
	return created, nil
}

// Delete ldap group.
// Returns 'ErrAlreadyNotExists' err if group already deleted.
// Returns 'ErrEntryNotDeleted' err when
// check after delete not returns 'ErrEntryNotFound'.
func (it *Groups) Delete(cn string) error {
	// Check if group already not exists
	g, err := it.GetByCn(cn)
	if err != nil {
		return err
	}
	if g == nil {
		return ErrAlreadyNotExists
	}

	dn := fmt.Sprintf("cn=%s,%s", cn, it.c.groupSearchBase)
	if err := it.c.cl.DeleteEntry(dn); err != nil {
		return fmt.Errorf("can`t delete group: %s", err.Error())
	}

	// Check that group realy deleted
	dg, err := it.GetByCn(cn)
	if err != nil {
		return fmt.Errorf("can`t check entry deletion: %s", err.Error())
	}
	if dg != nil {
		return fmt.Errorf("group not deleted")
	}
	return nil
}

// AddMembers Add users to a group members.
// Returns list of group members after update.
// Returns 'ErrEmptyMembersList' err when members list empty.
// Returns 'ErrNoNewMembersToAdd' err when group entry not updated.
func (it *Groups) AddMembers(groupCn string, membersUID []string) ([]string, error) {
	group, err := it.GetByCn(groupCn)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrEntryNotFound
	}
	if len(membersUID) == 0 {
		return nil, ErrEmptyMembersList
	}

	wg := &sync.WaitGroup{}
	ch := make(chan string, 500)
	for _, uid := range membersUID {
		wg.Add(1)
		go func(g *GroupEntry, member string, wg *sync.WaitGroup, ch chan<- string) {
			defer wg.Done()

			user, err := it.c.Users().GetByUidShort(member)
			if err != nil {
				log.Warnf("Can't get member '%s': %s", member, err.Error())
				return
			}
			if user == nil {
				log.Warnf("Member entry '%s' not found", member)
				return
			}
			if generigo.StringInSlice(user.ID, g.Members.Direct) {
				log.Warnf("Entry '%s' already member of a group", member)
				return
			}
			ch <- user.ID
		}(group, uid, wg, ch)
	}
	wg.Wait()
	close(ch)

	toadd := make([]string, 0, len(ch))
	for r := range ch {
		toadd = append(toadd, r)
	}
	if len(toadd) == 0 {
		return group.Members.Direct, ErrNoNewMembersToAdd
	}

	newmembers := make([]string, 0, len(toadd))
	newmembers = append(newmembers, group.Members.Direct...)
	newmembers = append(newmembers, toadd...)

	if err := it.c.cl.UpdateAttribute(group.DN, "memberUid", newmembers); err != nil {
		return nil, fmt.Errorf("can't update group members: %s", err.Error())
	}
	log.Info("Group members updated")

	updated, err := it.GetByCn(groupCn)
	if err != nil {
		return nil, fmt.Errorf("can't get group after update: %s", err.Error())
	}
	return updated.Members.Direct, nil
}

// DelMembers Delete users frmo a group members.
// Returns list of group members after update.
// Returns 'ErrEmptyMembersList' err when members list empty.
// Returns 'ErrNoNewMembersToDel' err when group entry not updated.
func (it *Groups) DelMembers(groupCn string, membersUID []string) ([]string, error) {
	g, err := it.GetByCn(groupCn)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, ErrEntryNotFound
	}
	if len(membersUID) == 0 {
		return nil, ErrEmptyMembersList
	}

	wg := &sync.WaitGroup{}
	ch := make(chan string, 500)
	for _, uid := range membersUID {
		wg.Add(1)
		go func(g *GroupEntry, member string, wg *sync.WaitGroup, ch chan<- string) {
			defer wg.Done()

			user, err := it.c.Users().GetByUidShort(member)
			if err != nil {
				log.Warnf("Can't get member '%s': %s", member, err.Error())
				return
			}
			if user == nil {
				log.Warnf("Member entry '%s' not found", member)
				return
			}
			if !generigo.StringInSlice(user.ID, g.Members.Direct) {
				log.Warnf("Entry '%s' already member of a group", member)
				return
			}
			ch <- user.ID
		}(g, uid, wg, ch)
	}
	wg.Wait()
	close(ch)

	todel := make([]string, 0, len(ch))
	for r := range ch {
		todel = append(todel, r)
	}
	if len(todel) == 0 {
		return nil, ErrNoNewMembersToDel
	}

	newmembers := make([]string, 0, len(g.Members.Direct))
	for _, currentMember := range g.Members.Direct {
		if !generigo.StringInSlice(currentMember, todel) {
			newmembers = append(newmembers, currentMember)
		}
	}

	if err := it.c.cl.UpdateAttribute(g.DN, "memberUid", newmembers); err != nil {
		return nil, fmt.Errorf("can't update group members: %s", err.Error())
	}
	log.Info("Group members updated")

	updated, err := it.GetByCn(groupCn)
	if err != nil {
		return nil, fmt.Errorf("can't get group after update: %s", err.Error())
	}
	return updated.Members.Direct, nil
}
