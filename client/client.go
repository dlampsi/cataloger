package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap/v3"
)

var (
	ErrTooManyEntriesFound = errors.New("too many entries found")
)

type Config struct {
	Host         string
	Port         int
	SSL          bool
	Insecure     bool
	BindDN       string
	BindPassword string
}

type Client struct {
	config *Config
}

// NewClient Initialize new Client and try to connect to server.
func New(config *Config) (*Client, error) {
	cl := Client{
		config: config,
	}
	if _, err := cl.Connect(); err != nil {
		return nil, err
	}
	return &cl, nil
}

// Connect Trying connect to ldap server.
func (cl *Client) Connect() (*ldap.Conn, error) {
	var (
		conn *ldap.Conn
		err  error
	)
	address := cl.config.Host + ":" + strconv.Itoa(cl.config.Port)
	if cl.config.SSL {
		conn, err = ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: cl.config.Insecure})
	} else {
		conn, err = ldap.Dial("tcp", address)
	}
	if err != nil {
		return nil, fmt.Errorf("dial error: %s", err.Error())
	}
	if err := conn.Bind(cl.config.BindDN, cl.config.BindPassword); err != nil {
		return nil, fmt.Errorf("bind error: %s", err.Error())
	}

	return conn, nil
}

// NewSearchRequest Returns ldap request object based in base search and filter.
// For use in Search methods.
func NewSearchRequest(base string, filter string) *ldap.SearchRequest {
	return &ldap.SearchRequest{
		BaseDN:       base,
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldap.NeverDerefAliases,
		SizeLimit:    0,
		TimeLimit:    0,
		TypesOnly:    false,
		Filter:       filter,
		Attributes:   []string{},
		Controls:     nil,
	}
}

// SearchEntries Perfroms search for ldap entries.
func (cl *Client) SearchEntries(r *ldap.SearchRequest) ([]*ldap.Entry, error) {
	con, err := cl.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	result, err := con.Search(r)
	if err != nil {
		return nil, err
	}
	return result.Entries, nil
}

// SearchEntry Perfrom search for single ldap entry.
// Returns nil if no entries found.
// Returns 'ErrTooManyEntriesFound' error if entries more that one.
func (cl *Client) SearchEntry(r *ldap.SearchRequest) (*ldap.Entry, error) {
	con, err := cl.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	result, err := con.Search(r)
	if err != nil {
		return nil, err
	}
	if len(result.Entries) > 1 {
		return nil, ErrTooManyEntriesFound
	}
	if len(result.Entries) < 1 {
		return nil, nil
	}
	return result.Entries[0], nil
}

// UpdateAttribute Performs update ldap entry attribure.
func (cl *Client) UpdateAttribute(dn string, attribute string, values []string) error {
	con, err := cl.Connect()
	if err != nil {
		return err
	}
	defer con.Close()
	mr := ldap.NewModifyRequest(dn, nil)
	mr.Replace(attribute, values)
	return con.Modify(mr)
}

// CreateEntry Creates new ldap entry.
func (cl *Client) CreateEntry(dn string, attributes map[string][]string) error {
	con, err := cl.Connect()
	if err != nil {
		return err
	}
	defer con.Close()
	ar := ldap.NewAddRequest(dn, nil)
	for k, v := range attributes {
		ar.Attribute(k, v)
	}
	return con.Add(ar)
}

// DeleteEntry Deletes ldap entry.
func (cl *Client) DeleteEntry(dn string) error {
	con, err := cl.Connect()
	if err != nil {
		return err
	}
	defer con.Close()
	dr := ldap.NewDelRequest(dn, nil)
	return con.Del(dr)
}
