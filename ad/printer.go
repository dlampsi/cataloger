package ad

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// Printer Prints user info to stdout.
func (it *Users) Printer(user *User) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "Parameter\tValue")
	fmt.Fprintln(w, "---------\t-----")
	fmt.Fprintln(w, "displayName\t", user.DisplayName)
	fmt.Fprintln(w, "sAMAccountName\t", user.ID)
	fmt.Fprintln(w, "dn\t", user.DN)
	fmt.Fprintln(w, "cn\t", user.CN)
	fmt.Fprintln(w, "mail\t", user.Mail)
	fmt.Fprintln(w)
	w.Flush()

	if len(user.Groups.All) != 0 {
		groups := ""
		for _, g := range user.Groups.All {
			groups = groups + fmt.Sprintf("%s\n", g)
		}
		fmt.Fprintln(w, "Groups")
		fmt.Fprintln(w, "------")
		fmt.Fprintln(w, groups)
		fmt.Fprintln(w)
		w.Flush()
	}

	// groupsDirect := ""
	// for _, g := range m.Groups.Direct {
	// 	groupsDirect = groupsDirect + fmt.Sprintf("%s\n", g)
	// }
	// fmt.Fprintln(w, "Groups (direct)")
	// fmt.Fprintln(w, "---------------")
	// fmt.Fprintln(w, groupsDirect)
	// fmt.Fprintln(w)
	// w.Flush()
}

// Printer Prints group info.
func (it *Groups) Printer(group *Group, printMembers bool) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "Parameter\tValue")
	fmt.Fprintln(w, "---------\t-----")
	fmt.Fprintln(w, "sAMAccountName\t", group.ID)
	fmt.Fprintln(w, "dn\t", group.DN)
	fmt.Fprintln(w, "cn\t", group.CN)
	fmt.Fprintln(w, "description\t", group.Description)
	fmt.Fprintln(w, "members_count\t", group.Members.Count)
	fmt.Fprintln(w)
	w.Flush()

	if printMembers {
		members := ""
		for _, member := range group.Members.All {
			members = members + fmt.Sprintf("%s\n", member)
		}
		fmt.Fprintln(w, "Members (sAMAccountName)")
		fmt.Fprintln(w, "------------------------")
		fmt.Fprintln(w, members)
		fmt.Fprintln(w)
		w.Flush()

		membersDN := ""
		for _, member := range group.Members.DirectDN {
			membersDN = membersDN + fmt.Sprintf("%s\n", member)
		}
		fmt.Fprintln(w, "Members (DN)")
		fmt.Fprintln(w, "------------")
		fmt.Fprintln(w, membersDN)
		fmt.Fprintln(w)
		w.Flush()
	}
}
