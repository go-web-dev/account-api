package repositories

type DBDriver interface {
	Close() error
}

type Accounts struct {
	db DBDriver
}

func NewAccounts(db DBDriver) Accounts {
	repo := Accounts{
		db,
	}
	return repo
}
