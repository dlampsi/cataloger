package ad

import (
	"cataloger/catalogs"
	"fmt"

	"github.com/dlampsi/generigo"
	"github.com/dlampsi/ldapconn"
	"gopkg.in/ldap.v3"
)

// Catalog AD catalog struct.
type Catalog struct {
	cl              *ldapconn.Client
	searchBase      string
	userSearchBase  string
	groupSearchBase string
}

// NewCatalog Initialize new AD catalog.
func NewCatalog(cfg *catalogs.Config) (*Catalog, error) {
	c := &Catalog{}

	c.searchBase = cfg.SearchBase

	c.userSearchBase = cfg.UserSearchBase
	if c.userSearchBase == "" {
		c.userSearchBase = c.searchBase
	}

	c.groupSearchBase = cfg.GroupSearchBase
	if c.groupSearchBase == "" {
		c.groupSearchBase = c.searchBase
	}

	// Create client
	cl, err := ldapconn.NewClient(&ldapconn.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		SSL:      cfg.SSL,
		Insecure: cfg.Insecure,
		BindDN:   cfg.BindDn,
		BindPass: cfg.BindPass,
	})
	if err != nil {
		return nil, fmt.Errorf("error create client: %s", err.Error())
	}
	c.cl = cl

	return c, nil
}

// EntryMembership Groups membership or group members structure.
type EntryMembership struct {
	Count    int      `json:"count"`
	All      []string `json:"all"`
	Direct   []string `json:"direct"`
	DirectDN []string `json:"direct_dn"`
}

// GetMembership Collect all AD groups
// where entry is member, including all subgroups.
func (c *Catalog) GetMembership(dn string, m *EntryMembership) (*EntryMembership, error) {
	if m == nil {
		m = &EntryMembership{}
	}
	filter := fmt.Sprintf(`(&(objectClass=group)(member=%v))`, ldap.EscapeFilter(dn))

	sr := ldapconn.CreateRequest(c.groupSearchBase, filter)
	ldapEntries, err := c.cl.SearchEntries(sr)
	if err != nil {
		return nil, fmt.Errorf("error get entry groups membership: " + err.Error())
	}

	// Add only first level groups DN
	if len(m.DirectDN) == 0 {
		for _, le := range ldapEntries {
			m.DirectDN = append(m.DirectDN, le.DN)
			m.Direct = append(m.Direct, le.GetAttributeValue("sAMAccountName"))
		}
	}

	for _, le := range ldapEntries {
		if !generigo.StringInSlice(le.GetAttributeValue("sAMAccountName"), m.All) {
			m.All = append(m.All, le.GetAttributeValue("sAMAccountName"))
			// Recourcive call entry group membership in other groups
			// m, _ = c.GetMembership(le.DN, m)
		}
	}

	return m, nil
}

// GetByDN Searches and returns ldap entry by DN.
// Returns 'ErrEntryNotFound' error when entry not founded.
func (c *Catalog) GetByDN(dn string) (*ldap.Entry, error) {
	filter := fmt.Sprintf("(distinguishedName=%s)", ldap.EscapeFilter(dn))
	sr := ldapconn.CreateRequest(c.searchBase, filter)
	return c.cl.SearchEntry(sr)
}

// CheckConn Check connection to catalog.
func (c *Catalog) CheckConn(cfg *catalogs.Config) error {
	_, err := ldapconn.NewClient(&ldapconn.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		SSL:      cfg.SSL,
		Insecure: cfg.Insecure,
		BindDN:   cfg.BindDn,
		BindPass: cfg.BindPass,
	})
	return err
}
