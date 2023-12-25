package ad

import (
	"cataloger/client"
	"fmt"
	"slices"

	"github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

// Membership AD entrie groups membership or group members structure.
type Membership struct {
	Count    int      `json:"count"`
	All      []string `json:"all"`
	Direct   []string `json:"direct"`
	DirectDN []string `json:"direct_dn"`
}

// GetByDN Searches and returns ldap entry by DN.
// Returns 'ErrEntryNotFound' error when entry not founded.
func (c *Catalog) GetByDN(dn string) (*ldap.Entry, error) {
	filter := fmt.Sprintf("(distinguishedName=%s)", ldap.EscapeFilter(dn))
	sr := client.NewSearchRequest(c.Attributes.SearchBase, filter)
	return c.cl.SearchEntry(sr)
}

// GetMembership Collect all AD groups where entry is member, including all subgroups.
func (c *Catalog) GetMembership(dn string, m *Membership) (*Membership, error) {
	if m == nil {
		m = &Membership{}
	}
	filter := fmt.Sprintf(`(&(objectClass=group)(member=%v))`, ldap.EscapeFilter(dn))
	log.WithFields(log.Fields{"filter": filter, "base": c.Attributes.SearchBase}).Debug("GetMembership")
	sr := client.NewSearchRequest(c.Attributes.SearchBase, filter)
	ldapEntries, err := c.cl.SearchEntries(sr)
	if err != nil {
		return nil, fmt.Errorf("can't get entry groups membership: " + err.Error())
	}
	// Add only first level groups DN
	if len(m.DirectDN) == 0 {
		for _, le := range ldapEntries {
			m.DirectDN = append(m.DirectDN, le.DN)
			m.Direct = append(m.Direct, le.GetAttributeValue("sAMAccountName"))
		}
	}
	for _, le := range ldapEntries {
		if !slices.Contains(m.All, le.GetAttributeValue("sAMAccountName")) {
			m.All = append(m.All, le.GetAttributeValue("sAMAccountName"))
			// Recourcive call entry group membership in other groups
			m, _ = c.GetMembership(le.DN, m)
		}
	}
	return m, nil
}
