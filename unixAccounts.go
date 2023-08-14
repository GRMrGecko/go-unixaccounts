package UNIXAccounts

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Base accounts structure and configuration.
type UNIXAccounts struct {
	Groups     []*UNIXGroup
	Users      []*UNIXUser
	PasswdPath string
	GroupPath  string
}

// Read the /etc/group and /etc/passwd files to parse information.
func NewUNIXAccounts() (*UNIXAccounts, error) {
	u := &UNIXAccounts{
		PasswdPath: "/etc/passwd",
		GroupPath:  "/etc/group",
	}
	err := u.Parse()
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Group data structure.
type UNIXGroup struct {
	Name  string
	ID    int
	Users []string
}

// User data structure.
type UNIXUser struct {
	Name     string
	ID       int
	GID      int
	FullName string
	HomeDir  string
	Shell    string
	Disabled bool
}

// Parse: Parse unix accounts and groups.
func (u *UNIXAccounts) Parse() error {
	// Remove any previously parsed users to re-parse.
	u.Groups = nil
	u.Users = nil

	// Open the group file.
	groupFile, err := os.Open(u.GroupPath)
	if err != nil {
		return err
	}
	defer groupFile.Close()

	scanner := bufio.NewScanner(groupFile)
	scanner.Split(bufio.ScanLines)
	lineCount := 0
	for scanner.Scan() {
		// Read a line.
		line := scanner.Text()

		// Ignore comments.
		if line[0] == '#' {
			continue
		}

		// Fields are separated with a :.
		fields := strings.Split(line, ":")

		// Groups should have 4 fields. Nothing more, nothing less.
		if len(fields) != 4 {
			return fmt.Errorf("unexpected field count in group file on line %d", lineCount)
		}

		// Parse information.
		group := new(UNIXGroup)
		group.Name = fields[0]
		group.ID, _ = strconv.Atoi(fields[2])
		group.Users = strings.Split(fields[3], ",")

		// Add group to array.
		u.Groups = append(u.Groups, group)

		// Increment line count.
		lineCount++
	}

	// Open the user file.
	passwdFile, err := os.Open(u.PasswdPath)
	if err != nil {
		return err
	}
	defer passwdFile.Close()

	scanner = bufio.NewScanner(passwdFile)
	scanner.Split(bufio.ScanLines)
	lineCount = 0
	for scanner.Scan() {
		// Read a line.
		line := scanner.Text()

		// Ignore comments.
		if line[0] == '#' {
			continue
		}

		// Fields are separated with a :.
		fields := strings.Split(line, ":")

		// Users have 7 fields. No more or less.
		if len(fields) != 7 {
			return fmt.Errorf("unexpected field count in passwd file on line %d", lineCount)
		}

		// Prase information.
		user := new(UNIXUser)
		user.Name = fields[0]
		user.ID, _ = strconv.Atoi(fields[2])
		user.GID, _ = strconv.Atoi(fields[3])
		user.FullName = fields[4]
		user.HomeDir = filepath.Clean(fields[5])
		user.Shell = fields[6]

		// A user is disabled if their shell is set to nologin or false. Users with no shell should also be disabled.
		user.Disabled = false
		if strings.Contains(user.Shell, "nologin") {
			user.Disabled = true
		}
		if strings.Contains(user.Shell, "false") {
			user.Disabled = true
		}
		if user.Shell == "" {
			user.Disabled = true
		}

		// Add user to array.
		u.Users = append(u.Users, user)

		// Increment line count.
		lineCount++
	}
	return nil
}

// Find user info for ID.
func (u *UNIXAccounts) UserWithID(id int) *UNIXUser {
	for _, user := range u.Users {
		if user.ID == id {
			return user
		}
	}
	return nil
}

// Find user info for name.
func (u *UNIXAccounts) UserWithName(name string) *UNIXUser {
	for _, user := range u.Users {
		if user.Name == name {
			return user
		}
	}
	return nil
}

// Find group info for ID.
func (u *UNIXAccounts) GroupWithID(id int) *UNIXGroup {
	for _, group := range u.Groups {
		if group.ID == id {
			return group
		}
	}
	return nil
}

// Find group info for name.
func (u *UNIXAccounts) GroupWithName(name string) *UNIXGroup {
	for _, group := range u.Groups {
		if group.Name == name {
			return group
		}
	}
	return nil
}

// Get all user accounts which are members of a group.
func (u *UNIXAccounts) UsersInGroup(group *UNIXGroup) []*UNIXUser {
	var users []*UNIXUser
	// Users with the Group ID set to the group's ID are a member.
	for _, user := range u.Users {
		if user.GID == group.ID {
			users = append(users, user)
		}
	}
	// Find user info for each member.
	for _, name := range group.Users {
		user := u.UserWithName(name)
		if user == nil {
			continue
		}
		// If the member was added previously, we do not want duplicates.
		alreadyExists := false
		for _, usr := range users {
			if usr == user {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			// The member is not a duplicate, so we add it to the array.
			users = append(users, user)
		}
	}
	return users
}

// List of groups a user is a member of.
func (u *UNIXAccounts) UserMemberOf(user *UNIXUser) []*UNIXGroup {
	var groups []*UNIXGroup

	// Look at each group and check if this user is a member.
	for _, group := range u.Groups {
		// If the GID of the user is this group, add it.
		if group.ID == user.GID {
			groups = append(groups, group)
			continue
		}

		// Check each user assigned to this group and add the group if the user matches.
		for _, thisUser := range group.Users {
			if thisUser == user.Name {
				// If the group was added previously, we do not want duplicates.
				alreadyExists := false
				for _, grp := range groups {
					if grp == group {
						alreadyExists = true
						break
					}
				}
				if !alreadyExists {
					// The group is not a duplicate, so we add it to the array.
					groups = append(groups, group)
				}
			}
		}
	}

	return groups
}
