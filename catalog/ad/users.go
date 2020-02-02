package ad

import (
	"fmt"

	"cataloger/client"

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

type User struct {
	ID          string     `json:"id"`
	DN          string     `json:"dn"`
	CN          string     `json:"cn"`
	Mail        string     `json:"mail"`
	DisplayName string     `json:"displayName"`
	Groups      Membership `json:"groups"`
}

// GetByAccountNameShort Get user info from catalog by user sAMAccountName.
// Returns nil structure if user not found.
func (it *Users) GetByFilter(filter string) ([]User, error) {
	log.WithFields(log.Fields{
		"filter": filter,
		"base":   it.c.Attributes.SearchBase,
	}).Debug("GetByFilter")
	sr := client.NewSearchRequest(it.c.Attributes.SearchBase, filter)
	entries, err := it.c.cl.SearchEntries(sr)
	if err != nil {
		return nil, fmt.Errorf("can't perform search: %s", err.Error())
	}
	if len(entries) == 0 {
		return nil, nil
	}
	users := []User{}
	for _, e := range entries {
		users = append(users, User{
			DN:          e.DN,
			ID:          e.GetAttributeValue("sAMAccountName"),
			CN:          e.GetAttributeValue("cn"),
			Mail:        e.GetAttributeValue("mail"),
			DisplayName: e.GetAttributeValue("displayName"),
		})
	}
	return users, nil
}

func (it *Users) GetGroups(user *User) (*User, error) {
	m, err := it.c.GetMembership(user.DN, nil)
	if err != nil {
		return nil, fmt.Errorf("can't get user membership: %s", err.Error())
	}
	user.Groups = Membership{
		All:      m.All,
		Direct:   m.Direct,
		DirectDN: m.DirectDN,
		Count:    len(m.All),
	}
	return user, nil
}
