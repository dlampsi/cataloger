package ad

import (
	"errors"
	"fmt"

	"github.com/dlampsi/ldapconn"
)

var (
	ErrNoNewMembersToAdd = errors.New("no new members to add")
	ErrNoNewMembersToDel = errors.New("no new members to delete")
	ErrEntryNotFound     = errors.New("entry not foiund")
	ErrEmptyMembersList  = errors.New("empty members list")
)

type Catalog struct {
	cl   *ldapconn.Client
	base string
}

func NewCatalog(ldapConf *ldapconn.Config, base string) (*Catalog, error) {
	c := &Catalog{}
	if base == "" {
		return nil, fmt.Errorf("empty search base provided")
	}
	c.base = base

	cl, err := ldapconn.NewClient(ldapConf)
	if err != nil {
		return nil, fmt.Errorf("can't init client: %s", err.Error())
	}
	c.cl = cl

	return c, nil
}
