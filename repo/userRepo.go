package repo

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}
type User struct {
	ID       int
	Login    string
	Password string
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(login, password string) error {
	_, err := r.db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", login, password)
	return err
}

func (r *UserRepository) GetUserByLogin(login string) (User, error) {
	var user User
	err := r.db.QueryRow("SELECT id, login, password FROM users WHERE login = $1", login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
