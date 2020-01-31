package ad

import (
	"fmt"

	"github.com/dlampsi/ldapconn"
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
