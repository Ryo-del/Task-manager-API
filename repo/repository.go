package repo

type User struct {
	ID       int
	Login    string
	Password string
}

type UserRepository interface {
	CreateUser(login, password string) error
	GetUserByLogin(login string) (User, error)
}
