package ldap

import (
	"cataloger/catalogs"
	"errors"
	"fmt"

	"github.com/dlampsi/ldapconn"
)

var (
	ErrNoNewMembersToAdd = errors.New("no new members to add")
	ErrNoNewMembersToDel = errors.New("no new members to delete")
	ErrEntryNotFound     = errors.New("entry not foiund")
	ErrEmptyMembersList  = errors.New("empty members list")
	ErrAlreadyExists     = errors.New("entry already exists in catalog")
	ErrAlreadyNotExists  = errors.New("entry already not exists in catalog")
)

// Catalog AD catalog struct.
type Catalog struct {
	cl              *ldapconn.Client
	searchBase      string
	userSearchBase  string
	groupSearchBase string
}

// NewCatalog Initialize new LDAP catalog.
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
	Direct   []string `json:"direct"`
	DirectDN []string `json:"direct_dn"`
}
