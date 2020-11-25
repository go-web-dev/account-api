package services

type AccountsRepository interface {
}

func NewAuth(repo AccountsRepository) Auth {
	service := Auth{
		repo: repo,
	}
	return service
}

// Auth represents the Auth service
type Auth struct {
	repo AccountsRepository
}

func (s Auth) Login() {
}

func (s Auth) Signup() {
}

func (s Auth) Logout() {
}
