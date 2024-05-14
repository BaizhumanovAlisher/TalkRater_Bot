package databases

// Application can save admins usernames in local map
// because it will contain about 1-5 nicknames
// Also admins are defined by the config file
// and can only be changed when the application is restarted

type AdminDB struct {
	admins map[string]struct{}
}

func NewAdminDB(adminsStr []string) *AdminDB {
	db := &AdminDB{
		admins: make(map[string]struct{}),
	}

	for _, admin := range adminsStr {
		db.admins[admin] = struct{}{}
	}

	return db
}

func (db *AdminDB) IsAdmin(username string) bool {
	_, ok := db.admins[username]
	return ok
}
