package ad

import (
	"fmt"

	"github.com/dlampsi/ldapconn"
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

// GetByAccountNameShort Get user info from catalog by user sAMAccountName.
// Returns nil structure if user not found.
func (it *Users) GetByAccountNameShort(sAMAccountName string) (*UserEntry, error) {
	filter := fmt.Sprintf("(sAMAccountName=%s)", sAMAccountName)
	sr := ldapconn.CreateRequest(it.c.searchBase, filter)
	entry, err := it.c.cl.SearchEntry(sr)
	if err != nil {
		return nil, fmt.Errorf("error get: " + err.Error())
	}
	if entry == nil {
		return nil, nil
	}
	u := &UserEntry{
		DN:          entry.DN,
		ID:          entry.GetAttributeValue("sAMAccountName"),
		CN:          entry.GetAttributeValue("cn"),
		Mail:        entry.GetAttributeValue("mail"),
		DisplayName: entry.GetAttributeValue("displayName"),
	}

	return u, nil
}

// GetByAccountName Returns AD user info. Searches by sAMAccountName.
func (it *Users) GetByAccountName(sAMAccountName string) (*UserEntry, error) {
	user, err := it.GetByAccountNameShort(sAMAccountName)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	membsh, err := it.c.GetMembership(user.DN, nil)
	if err != nil {
		return nil, err
	}
	user.Groups = EntryMembership{
		All:      membsh.All,
		Direct:   membsh.Direct,
		DirectDN: membsh.DirectDN,
		Count:    len(membsh.All),
	}

	return user, nil
}
