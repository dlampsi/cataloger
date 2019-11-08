package ad

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

// Groups Entry point to AD catalog groups methods.
func (c *Catalog) Groups() *Groups {
	return &Groups{c: c}
}

// GroupEntry Grouper group entry struct.
type GroupEntry struct {
	ID            string          `json:"id"`
	DN            string          `json:"dn"`
	CN            string          `json:"cn"`
	Description   string          `json:"description"`
	Members       EntryMembership `json:"members"`
	MembersGroups []string
}

// GetShort Returns short AD group info without all group members search.
// Created for fast requests when members not needed.
// Returns nil structure if group not found.
func (it *Groups) GetByAccountNameShort(sAMAccountName string) (*GroupEntry, error) {
	filter := fmt.Sprintf("(&(objectClass=group)(sAMAccountName=%s))", sAMAccountName)
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
		ID:          ldapEntry.GetAttributeValue("sAMAccountName"),
		CN:          ldapEntry.GetAttributeValue("cn"),
		Description: ldapEntry.GetAttributeValue("description"),
		Members: EntryMembership{
			DirectDN: ldapEntry.GetAttributeValues("member"),
		},
	}, nil
}

// Get Returns AD group info.
func (it *Groups) GetByAccountName(sAMAccountName string) (*GroupEntry, error) {
	group, err := it.GetByAccountNameShort(sAMAccountName)
	if err != nil {
		return nil, err
	}

	group.Members.Direct = make([]string, 0, len(group.Members.DirectDN))
	group.Members.All = []string{}

	wg := &sync.WaitGroup{}
	directCh := make(chan string, len(group.Members.DirectDN))
	allCh := make(chan string, 5000)

	// Using quota chanel because ldap connection may be refused
	// because of too many connection opened by go routines.
	quotaCh := make(chan struct{}, 10)

	// Loop over all group members DNs and fill members sAMAccountName lists
	for _, memberDN := range group.Members.DirectDN {
		wg.Add(1)
		go it.procGroupMembers(memberDN, wg, directCh, allCh, quotaCh)
	}
	wg.Wait()

	close(directCh)
	close(allCh)

	for id := range directCh {
		group.Members.Direct = append(group.Members.Direct, id)
	}
	for id := range allCh {
		if !generigo.StringInSlice(id, group.Members.All) {
			group.Members.All = append(group.Members.All, id)
		}
	}

	group.Members.Count = len(group.Members.All)

	return group, nil
}

// Process group members.
// Get member by DN, determine if member is group and recoursive search members.
// Use limitation on quotaCh.
func (it *Groups) procGroupMembers(memberDN string, wg *sync.WaitGroup, directCh chan<- string, allCh chan<- string, quotaCh chan struct{}) {
	quotaCh <- struct{}{} // Take a slot in the quota channel
	defer wg.Done()

	ldapEntry, err := it.c.GetByDN(memberDN)
	if err != nil {
		fmt.Printf("Can't get group member by dn: %s\n", err.Error())
		return
	}
	memberID := ldapEntry.GetAttributeValue("sAMAccountName")

	directCh <- memberID

	// If group member is a group, recoursive search members.
	if generigo.StringInSlice("group", ldapEntry.GetAttributeValues("objectClass")) {
		subgroup, err := it.GetByAccountName(memberID)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, sm := range subgroup.Members.All {
			allCh <- sm
		}
	} else {
		allCh <- memberID
	}

	<-quotaCh // Release slot in the quota channel
}

// AddMembersByAccountName Add group members.
// Proc all new members, add to chanel only if member exists and he's not a member
func (it *Groups) AddMembersByAccountName(sAMAccountName string, members []string) error {
	log.Debugf("Add '%s' members: %v", sAMAccountName, members)
	group, err := it.GetByAccountName(sAMAccountName)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrEntryNotFound
	}
	log.Debugf("Group '%s' found", sAMAccountName)
	if len(members) == 0 {
		return ErrEmptyMembersList
	}

	wg := &sync.WaitGroup{}
	ch := make(chan string, 500)
	for _, member := range members {
		wg.Add(1)
		go func(group *GroupEntry, member string, wg *sync.WaitGroup, ch chan<- string) {
			defer wg.Done()
			log.Debugf("Processing member '%s'", member)
			user, err := it.c.Users().GetByAccountNameShort(member)
			if err != nil {
				log.Warnf("Can't get member '%s': %s", member, err.Error())
				return
			}
			if user == nil {
				log.Warnf("Member entry '%s' not found", member)
				return
			}
			if generigo.StringInSlice(user.ID, group.Members.All) {
				log.Warnf("Entry '%s' already member of a group", member)
				return
			}
			log.Debugf("New member '%s' will be added", member)
			ch <- user.DN
		}(group, member, wg, ch)
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

	newMembersDN := make([]string, 0, len(group.Members.DirectDN)+len(toadd))
	newMembersDN = append(newMembersDN, group.Members.DirectDN...)
	newMembersDN = append(newMembersDN, toadd...)

	// Update
	if err := it.c.cl.UpdateAttribute(group.DN, "member", newMembersDN); err != nil {
		return err
	}

	log.Infof("Added %d group members", len(toadd))

	return nil
}

// DelMembersByAccountName Delete user from a group members.
// Delete from group only direct group members.
func (it *Groups) DelMembersByAccountName(sAMAccountName string, members []string) error {
	log.Debugf("Del '%s' members: %v", sAMAccountName, members)
	group, err := it.GetByAccountName(sAMAccountName)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrEntryNotFound
	}
	log.Debugf("Group '%s' found", sAMAccountName)
	if len(members) == 0 {
		return ErrEmptyMembersList
	}

	wg := &sync.WaitGroup{}
	ch := make(chan string, 500)
	for _, member := range members {
		wg.Add(1)
		go func(group *GroupEntry, member string, wg *sync.WaitGroup, ch chan<- string) {
			defer wg.Done()
			log.Debugf("Processing member '%s'", member)
			user, err := it.c.Users().GetByAccountNameShort(member)
			if err != nil {
				log.Warnf("Can't get member '%s': %s", member, err.Error())
				return
			}
			if user == nil {
				log.Warnf("Member entry '%s' not found", member)
				return
			}
			if !generigo.StringInSlice(user.ID, group.Members.All) {
				log.Warnf("Entry '%s' already not a member of a group or a subgroup", user.ID)
				return
			}
			log.Debugf("Member '%s' will be deleted from group", member)
			ch <- user.DN
		}(group, member, wg, ch)
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

	// Delete from group only current group members
	newMembersDN := make([]string, 0, len(group.Members.DirectDN))
	for _, curMember := range group.Members.DirectDN {
		if !generigo.StringInSlice(curMember, todell) {
			newMembersDN = append(newMembersDN, curMember)
		}
	}

	// Update
	if err := it.c.cl.UpdateAttribute(group.DN, "member", newMembersDN); err != nil {
		return fmt.Errorf("error update group: %s", err.Error())
	}

	log.Infof("Removed %d group members", len(todell))

	return nil
}
