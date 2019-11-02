package ad

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"

	"github.com/dlampsi/generigo"
	"github.com/dlampsi/ldapconn"
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

// Print Prints group info.
func (it *Groups) Print(m *GroupEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "Parameter\tValue")
	fmt.Fprintln(w, "---------\t-----")
	fmt.Fprintln(w, "sAMAccountName\t", m.ID)
	fmt.Fprintln(w, "dn\t", m.DN)
	fmt.Fprintln(w, "cn\t", m.CN)
	fmt.Fprintln(w, "description\t", m.Description)
	fmt.Fprintln(w, "members_count\t", len(m.Members.All))
	fmt.Fprintln(w)
	w.Flush()

	members := ""
	for _, member := range m.Members.All {
		members = members + fmt.Sprintf("%s\n", member)
	}
	fmt.Fprintln(w, "Members")
	fmt.Fprintln(w, "------")
	fmt.Fprintln(w, members)
	fmt.Fprintln(w)
	w.Flush()

	membersDirect := ""
	for _, member := range m.Members.Direct {
		membersDirect = membersDirect + fmt.Sprintf("%s\n", member)
	}
	fmt.Fprintln(w, "Members (direct)")
	fmt.Fprintln(w, "----------------")
	fmt.Fprintln(w, membersDirect)
	fmt.Fprintln(w)
	w.Flush()
}

// PrintMembers Prints group members list only.
func (it *Groups) PrintMembers(m *GroupEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)
	members := ""
	for _, member := range m.Members.All {
		members = members + fmt.Sprintf("%s\n", member)
	}
	fmt.Fprintln(w, "Members")
	fmt.Fprintln(w, "------")
	fmt.Fprintln(w, members)
	fmt.Fprintln(w)
	w.Flush()
}

// PrintMembersDirect Prints group members list only.
func (it *Groups) PrintMembersDirect(m *GroupEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)
	membersDirect := ""
	for _, member := range m.Members.Direct {
		membersDirect = membersDirect + fmt.Sprintf("%s\n", member)
	}
	fmt.Fprintln(w, "Members (direct)")
	fmt.Fprintln(w, "----------------")
	fmt.Fprintln(w, membersDirect)
	fmt.Fprintln(w)
	w.Flush()
}

// // AddMembers Add user to group members.
// // Param groupID is sAMAccountName attribute.
// // Param membersIDs is a list of sAMAccountName attributes.
// func (it *Groups) AddMembersByAccountName(sAMAccountName string, membersIDs []string) error {
// 	group, err := it.GetByAccountName(sAMAccountName)
// 	if err != nil {
// 		return err
// 	}
// 	if len(membersIDs) == 0 {
// 		return errors.New("empty method input data")
// 	}

// 	wg := &sync.WaitGroup{}
// 	ch := make(chan string, 500)
// 	for _, memberID := range membersIDs {
// 		wg.Add(1)
// 		go func(group *GroupEntry, memberID string, wg *sync.WaitGroup, ch chan<- string) {
// 			defer wg.Done()
// 			u, err := it.c.Users().GetByAccountNameShort(memberID)
// 			if err != nil {
// 				rlog.Error(err)
// 				return
// 			}
// 			if u == nil {
// 				fmt.Printf("Member '%s' not found", memberID)
// 			}
// 			if generigo.StringInSlice(u.ID, group.Members.All) {
// 				rlog.Warn("Already member of a group")
// 				return
// 			}
// 			ch <- u.DN
// 		}(group, memberID, wg, ch)
// 	}

// 	wg.Wait()
// 	close(ch)

// 	toadd := make([]string, 0, len(ch))
// 	for r := range ch {
// 		toadd = append(toadd, r)
// 	}

// 	if len(toadd) == 0 {
// 		return ErrEntryNotUpdated
// 	}

// 	newMembersDN := make([]string, 0, len(group.Members.DirectDN)+len(toadd))
// 	newMembersDN = append(newMembersDN, group.Members.DirectDN...)
// 	newMembersDN = append(newMembersDN, toadd...)

// 	// Update
// 	if err := it.c.cl.UpdateAttribute(group.DN, "member", newMembersDN); err != nil {
// 		flog.Errorf("Error update ldap entry: %s", err.Error())
// 		return err
// 	}
// 	flog.Info("Group members updated")

// 	return nil
// }
