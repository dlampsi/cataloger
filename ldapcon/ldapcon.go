package ldapcon

import (
	"crypto/tls"
	"fmt"
	"strconv"

	"gopkg.in/ldap.v3"
)

// Config Ldap server connection parameters.
type Config struct {
	Host     string
	Port     int
	SSL      bool
	Insecure bool
	BindDN   string
	BindPass string
}

// Client struct.
type Client struct {
	config *Config
}

// NewClient Initialize new ldap client.
func NewClient(cfg *Config) (*Client, error) {
	cl := &Client{
		config: cfg,
	}

	if _, err := cl.Connect(); err != nil {
		return nil, err
	}

	return cl, nil
}

// Connect Trying connect ldap server.
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
	if err := conn.Bind(cl.config.BindDN, cl.config.BindPass); err != nil {
		return nil, fmt.Errorf("bind error: %s", err.Error())
	}
	return conn, nil
}

// SearchEntries Perfrom search for ldap entries.
func (cl *Client) SearchEntries(sr *ldap.SearchRequest) ([]*ldap.Entry, error) {
	con, err := cl.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()

	result, err := con.Search(sr)
	if err != nil {
		return nil, err
	}

	return result.Entries, nil
}

// SearchEntry Perfrom search for single ldap entry.
func (cl *Client) SearchEntry(sr *ldap.SearchRequest) (*ldap.Entry, error) {
	con, err := cl.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()

	result, err := con.Search(sr)
	if err != nil {
		return nil, err
	}

	if len(result.Entries) > 1 {
		return nil, fmt.Errorf("to many ldap entries finded (" + strconv.Itoa(len(result.Entries)) + ")")
	}

	if len(result.Entries) < 1 {
		return nil, nil
	}

	return result.Entries[0], nil
}

// UpdateAttribute Perform update ldap entry attribure.
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

// DeleteEntry Delete ldap entry.
func (cl *Client) DeleteEntry(dn string) error {
	con, err := cl.Connect()
	if err != nil {
		return err
	}
	defer con.Close()

	dr := ldap.NewDelRequest(dn, nil)

	return con.Del(dr)
}

// CreateRequest Return ldap request object based in base search and filter.
// For use in Search methods.
func CreateRequest(base string, filter string) *ldap.SearchRequest {
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
