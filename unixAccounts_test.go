package UNIXAccounts

import "testing"

func TestAccounts(t *testing.T) {
	u := &UNIXAccounts{
		PasswdPath: "test/passwd",
		GroupPath:  "test/group",
	}
	err := u.Parse()
	if err != nil {
		t.Errorf("error parsing: %s", err)
	}

	user := u.UserWithID(2)
	if user == nil || user.Name != "daemon" {
		t.Error("unexpected user found by id")
	}

	user = u.UserWithName("test")
	if user == nil || user.ID != 1000 {
		t.Error("unexpected user found by name")
	}

	group := u.GroupWithID(1)
	if group == nil || group.Name != "bin" {
		t.Error("unexpected group found by id")
	}

	group = u.GroupWithName("cdrom")
	if group == nil || group.ID != 11 {
		t.Error("unexpected group found by name")
	}

	users := u.UsersInGroup(group)
	if len(users) != 2 {
		t.Error("unexpected user count found")
	}
	for _, usr := range users {
		if usr.Name != "root" && usr.Name != "test" {
			t.Errorf("found unexpected user in group: %s", usr.Name)
		}
	}

	groups := u.UserMemberOf(user)
	if len(groups) != 2 {
		t.Error("unexpected group count found")
	}
	for _, grp := range groups {
		if grp.Name != "test" && grp.Name != "cdrom" {
			t.Errorf("found unexpected group found by user: %s", grp.Name)
		}
	}

	u = &UNIXAccounts{
		PasswdPath: "test/invalid-passwd",
		GroupPath:  "test/group",
	}
	err = u.Parse()
	if err == nil {
		t.Error("expected parse to fail, but it succeeded.")
	}

	u = &UNIXAccounts{
		PasswdPath: "test/passwd",
		GroupPath:  "test/invalid-group",
	}
	err = u.Parse()
	if err == nil {
		t.Error("expected parse to fail, but it succeeded.")
	}
}
