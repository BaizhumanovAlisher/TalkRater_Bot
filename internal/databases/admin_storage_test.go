package databases

import "testing"

func TestAdminDB(t *testing.T) {
	tests := []struct {
		name          string
		initialAdmins []string
		checkUsers    []struct {
			username string
			isAdmin  bool
		}
	}{
		{
			name:          "Single admin",
			initialAdmins: []string{"admin1"},
			checkUsers: []struct {
				username string
				isAdmin  bool
			}{
				{username: "admin1", isAdmin: true},
				{username: "user1", isAdmin: false},
			},
		},
		{
			name:          "Multiple admins",
			initialAdmins: []string{"admin1", "admin2"},
			checkUsers: []struct {
				username string
				isAdmin  bool
			}{
				{username: "admin1", isAdmin: true},
				{username: "admin2", isAdmin: true},
				{username: "user1", isAdmin: false},
			},
		},
		{
			name:          "No admins",
			initialAdmins: []string{},
			checkUsers: []struct {
				username string
				isAdmin  bool
			}{
				{username: "admin1", isAdmin: false},
				{username: "user1", isAdmin: false},
			},
		},
		{
			name:          "Non-existing admin",
			initialAdmins: []string{"admin1"},
			checkUsers: []struct {
				username string
				isAdmin  bool
			}{
				{username: "admin2", isAdmin: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewAdminDB(tt.initialAdmins)
			for _, check := range tt.checkUsers {
				result := db.IsAdmin(check.username)
				if result != check.isAdmin {
					t.Errorf("Expected IsAdmin(%s) to be %v, got %v", check.username, check.isAdmin, result)
				}
			}
		})
	}
}
