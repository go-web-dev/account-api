package repositories

// DBDriver represents the database driver
type DBDriver interface {
	Close() error
}

// Accounts represents the accounts repository
type Accounts struct {
	db DBDriver
}

// NewAccounts creates a new Accounts repository
func NewAccounts(db DBDriver) Accounts {
	repo := Accounts{
		db,
	}
	return repo
}
