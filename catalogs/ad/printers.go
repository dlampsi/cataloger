package ad

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
	fmt.Fprintln(w, "displayName\t", m.DisplayName)
	fmt.Fprintln(w, "sAMAccountName\t", m.ID)
	fmt.Fprintln(w, "dn\t", m.DN)
	fmt.Fprintln(w, "cn\t", m.CN)
	fmt.Fprintln(w, "mail\t", m.Mail)
	// fmt.Fprintln(w, "groups_count\t", len(m.Groups.All))
	fmt.Fprintln(w)
	w.Flush()

	groups := ""
	for _, g := range m.Groups.All {
		groups = groups + fmt.Sprintf("%s\n", g)
	}
	fmt.Fprintln(w, "Groups")
	fmt.Fprintln(w, "------")
	fmt.Fprintln(w, groups)
	fmt.Fprintln(w)
	w.Flush()

	groupsDirect := ""
	for _, g := range m.Groups.Direct {
		groupsDirect = groupsDirect + fmt.Sprintf("%s\n", g)
	}
	fmt.Fprintln(w, "Groups (direct)")
	fmt.Fprintln(w, "---------------")
	fmt.Fprintln(w, groupsDirect)
	fmt.Fprintln(w)
	w.Flush()
}

// Print user groups list.
func (it *Users) PrintGroups(m *UserEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)
	groups := ""
	for _, g := range m.Groups.All {
		groups = groups + fmt.Sprintf("%s\n", g)
	}
	fmt.Fprintln(w, "Groups")
	fmt.Fprintln(w, "------")
	fmt.Fprintln(w, groups)
	fmt.Fprintln(w)
	w.Flush()
}

// Print only direct user groups list.
func (it *Users) PrintGroupsDirect(m *UserEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)
	groupsDirect := ""
	for _, g := range m.Groups.Direct {
		groupsDirect = groupsDirect + fmt.Sprintf("%s\n", g)
	}
	fmt.Fprintln(w, "Groups (direct)")
	fmt.Fprintln(w, "---------------")
	fmt.Fprintln(w, groupsDirect)
	fmt.Fprintln(w)
	w.Flush()
}

// Print Prints group info.
func (it *Groups) Print(m *GroupEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "Parameter\tValue")
	fmt.Fprintln(w, "---------\t-----")
	fmt.Fprintln(w, "sAMAccountName\t", m.ID)
	fmt.Fprintln(w, "dn\t", m.DN)
	fmt.Fprintln(w, "cn\t", m.CN)
	fmt.Fprintln(w, "description\t", m.Description)
	fmt.Fprintln(w, "members_count\t", len(m.Members.All))
	fmt.Fprintln(w)
	w.Flush()

	members := ""
	for _, member := range m.Members.All {
		members = members + fmt.Sprintf("%s\n", member)
	}
	fmt.Fprintln(w, "Members")
	fmt.Fprintln(w, "------")
	fmt.Fprintln(w, members)
	fmt.Fprintln(w)
	w.Flush()

	membersDirect := ""
	for _, member := range m.Members.Direct {
		membersDirect = membersDirect + fmt.Sprintf("%s\n", member)
	}
	fmt.Fprintln(w, "Members (direct)")
	fmt.Fprintln(w, "----------------")
	fmt.Fprintln(w, membersDirect)
	fmt.Fprintln(w)
	w.Flush()
}

// PrintMembers Prints group members list only.
func (it *Groups) PrintMembers(m *GroupEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)
	members := ""
	for _, member := range m.Members.All {
		members = members + fmt.Sprintf("%s\n", member)
	}
	fmt.Fprintln(w, "Members")
	fmt.Fprintln(w, "------")
	fmt.Fprintln(w, members)
	fmt.Fprintln(w)
	w.Flush()
}

// PrintMembersDirect Prints group members list only.
func (it *Groups) PrintMembersDirect(m *GroupEntry) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)
	membersDirect := ""
	for _, member := range m.Members.Direct {
		membersDirect = membersDirect + fmt.Sprintf("%s\n", member)
	}
	fmt.Fprintln(w, "Members (direct)")
	fmt.Fprintln(w, "----------------")
	fmt.Fprintln(w, membersDirect)
	fmt.Fprintln(w)
	w.Flush()
}
