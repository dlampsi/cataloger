package ad

import (
	"fmt"

	"github.com/dlampsi/ldapconn"
	log "github.com/sirupsen/logrus"
)

// Users Catalog users struct.
type Users struct {
	c *Catalog
}

// Users Entry point to AD catalog users methods.
func (c *Catalog) Users() *Users {
	return &Users{c: c}
}

type UserEntry struct {
	ID          string          `json:"id"`
	DN          string          `json:"dn"`
	CN          string          `json:"cn"`
	Mail        string          `json:"mail"`
	DisplayName string          `json:"displayName"`
	Groups      EntryMembership `json:"groups"`
}

// Get Search user info in catalog by specified ldap filter.
// Returns nil structure if user not found.
func (it *Users) Get(filter string) (*UserEntry, error) {
	u, err := it.GetShort(filter)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	membsh, err := it.c.GetMembership(u.DN, nil)
	if err != nil {
		return nil, err
	}
	u.Groups = EntryMembership{
		All:      membsh.All,
		Direct:   membsh.Direct,
		DirectDN: membsh.DirectDN,
		Count:    len(membsh.All),
	}
	return u, nil
}

// GetShort Search user info in catalog by specified ldap filter.
// Without user group membership.
// Returns nil structure if user not found.
func (it *Users) GetShort(filter string) (*UserEntry, error) {
	sr := ldapconn.CreateRequest(it.c.searchBase, filter)
	entry, err := it.c.cl.SearchEntry(sr)
	if err != nil {
		return nil, fmt.Errorf("error get: " + err.Error())
	}
	if entry == nil {
		log.Debugf("No entries found by filter '%s'", filter)
		return nil, nil
	}
	return &UserEntry{
		DN:          entry.DN,
		ID:          entry.GetAttributeValue("sAMAccountName"),
		CN:          entry.GetAttributeValue("cn"),
		Mail:        entry.GetAttributeValue("mail"),
		DisplayName: entry.GetAttributeValue("displayName"),
	}, nil
}

// GetByAccountName Get user info from catalog by user sAMAccountName.
// Returns nil structure if user not found.
func (it *Users) GetByAccountName(sAMAccountName string) (*UserEntry, error) {
	filter := fmt.Sprintf("(sAMAccountName=%s)", sAMAccountName)
	return it.Get(filter)
}

// GetByAccountNameShort Get user info from catalog by user sAMAccountName.
// Without user group membership.
// Returns nil structure if user not found.
func (it *Users) GetByAccountNameShort(sAMAccountName string) (*UserEntry, error) {
	filter := fmt.Sprintf("(sAMAccountName=%s)", sAMAccountName)
	return it.GetShort(filter)
}
