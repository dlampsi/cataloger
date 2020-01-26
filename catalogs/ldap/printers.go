package ldap

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// Print user info.
func (it *Users) Print(m *UserEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "Parameter\tValue")
	fmt.Fprintln(w, "---------\t-----")
	fmt.Fprintln(w, "uid\t", m.ID)
	fmt.Fprintln(w, "dn\t", m.DN)
	fmt.Fprintln(w, "cn\t", m.CN)
	fmt.Fprintln(w, "mail\t", m.Mail)
	// fmt.Fprintln(w, "groups_count\t", len(m.Groups.All))
	fmt.Fprintln(w)
	w.Flush()

	groups := ""
	for _, g := range m.Groups.Direct {
		groups = groups + fmt.Sprintf("%s\n", g)
	}
	fmt.Fprintln(w, "Groups")
	fmt.Fprintln(w, "------")
	fmt.Fprintln(w, groups)
	fmt.Fprintln(w)
	w.Flush()
}

// Print user groups list.
func (it *Users) PrintGroups(m *UserEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)
	groups := ""
	for _, g := range m.Groups.Direct {
		groups = groups + fmt.Sprintf("%s\n", g)
	}
	fmt.Fprintln(w, "Groups")
	fmt.Fprintln(w, "------")
	fmt.Fprintln(w, groups)
	fmt.Fprintln(w)
	w.Flush()
}

// Print Prints group info.
func (it *Groups) Print(m *GroupEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "Parameter\tValue")
	fmt.Fprintln(w, "---------\t-----")
	fmt.Fprintln(w, "uid\t", m.ID)
	fmt.Fprintln(w, "dn\t", m.DN)
	fmt.Fprintln(w, "cn\t", m.CN)
	fmt.Fprintln(w, "description\t", m.Description)
	fmt.Fprintln(w, "members_count\t", len(m.Members.Direct))
	fmt.Fprintln(w)
	w.Flush()

	members := ""
	for _, member := range m.Members.Direct {
		members = members + fmt.Sprintf("%s\n", member)
	}
	fmt.Fprintln(w, "Members")
	fmt.Fprintln(w, "------")
	fmt.Fprintln(w, members)
	fmt.Fprintln(w)
	w.Flush()
}

// PrintMembers Prints group members list only.
func (it *Groups) PrintMembers(m *GroupEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)
	members := ""
	for _, member := range m.Members.Direct {
		members = members + fmt.Sprintf("%s\n", member)
	}
	fmt.Fprintln(w, "Members")
	fmt.Fprintln(w, "------")
	fmt.Fprintln(w, members)
	fmt.Fprintln(w)
	w.Flush()
}
