package ldap

import (
	"fmt"

	"github.com/dlampsi/ldapconn"
)

// Users Catalog users struct.
type Users struct {
	c *Catalog
}

// Users Entry point to LDAP catalog users methods.
func (c *Catalog) Users() *Users {
	return &Users{c: c}
}

type UserEntry struct {
	ID        string          `json:"id"`
	DN        string          `json:"dn"`
	CN        string          `json:"cn"`
	Mail      string          `json:"mail"`
	Groups    EntryMembership `json:"groups"`
	Netgroups EntryMembership `json:"netgroups"`
}

// GetByUidShort Get short user info (withour groups) from catalog by user uid.
// Returns nil structure if user not found.
func (it *Users) GetByUidShort(uid string) (*UserEntry, error) {
	filter := fmt.Sprintf("(uid=%s)", uid)
	sr := ldapconn.CreateRequest(it.c.searchBase, filter)
	entry, err := it.c.cl.SearchEntry(sr)
	if err != nil {
		return nil, fmt.Errorf("error get: " + err.Error())
	}
	if entry == nil {
		return nil, nil
	}
	u := &UserEntry{
		DN:   entry.DN,
		ID:   entry.GetAttributeValue("uid"),
		CN:   entry.GetAttributeValue("cn"),
		Mail: entry.GetAttributeValue("mail"),
	}

	return u, nil
}

// GetByUid Get user info from catalog by user uid.
// Returns nil structure if user not found.
func (it *Users) GetByUid(uid string) (*UserEntry, error) {
	user, err := it.GetByUidShort(uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	ug, err := it.getGroups(uid)
	if err != nil {
		return nil, err
	}
	user.Groups = EntryMembership{
		Count:    ug.Count,
		Direct:   ug.Direct,
		DirectDN: ug.DirectDN,
	}

	return user, nil
}

// Return user groups struct.
func (it *Users) getGroups(uid string) (*EntryMembership, error) {
	m := &EntryMembership{}
	filter := fmt.Sprintf("(&(objectClass=posixGroup)(memberUid=%s))", uid)
	srG := ldapconn.CreateRequest(it.c.groupSearchBase, filter)
	ldapEntries, err := it.c.cl.SearchEntries(srG)
	if err != nil {
		return nil, fmt.Errorf("error get entry: %s", err.Error())
	}
	for _, le := range ldapEntries {
		m.DirectDN = append(m.DirectDN, le.DN)
		m.Direct = append(m.Direct, le.GetAttributeValue("cn"))
	}
	m.Count = len(m.Direct)
	return m, nil
}
